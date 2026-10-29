package processor

import (
	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
)

// AnalyzedPacket contains the extracted key information from a raw network packet.
type AnalyzedPacket struct {
	Timestamp int64
	SourceIP  string
	DestIP    string
	SourcePort int
	DestPort   int
	Protocol   string
	PayloadLen int
}

// Analyzer processes raw packets and extracts relevant network information.
type Analyzer struct {
	filterPortMin int
	filterPortMax int
}

// NewAnalyzer creates a new Analyzer with specified port range filters for "game-like" traffic.
func NewAnalyzer(portMin, portMax int) *Analyzer {
	return &Analyzer{
		filterPortMin: portMin,
		filterPortMax: portMax,
	}
}

// Process consumes packets from the reader channel and returns a channel of analyzed packets.
func (a *Analyzer) Process(packets <-chan Packet) <-chan AnalyzedPacket {
	analyzedChan := make(chan AnalyzedPacket, 1000)

	go func() {
		defer close(analyzedChan)
		for p := range packets {
			analyzed, ok := a.analyzePacket(p)
			if ok {
				analyzedChan <- analyzed
			}
		}
	}()

	return analyzedChan
}

func (a *Analyzer) analyzePacket(p Packet) (AnalyzedPacket, bool) {
	// Use gopacket to decode the raw payload.
	// Since we are dealing with raw payloads from the sniffer,
	// we assume they start with the Ethernet layer.
	packet := gopacket.NewPacket(p.Payload, layers.LayerTypeEthernet, gopacket.Default)

	// Extract Network Layer (IPv4 or IPv6)
	var srcIP, dstIP string
	var protocol string

	if ip4Layer := packet.Layer(layers.LayerTypeIPv4); ip4Layer != nil {
		ip4 := ip4Layer.(*layers.IPv4)
		srcIP = ip4.SrcIP.String()
		dstIP = ip4.DstIP.String()
		protocol = ip4.Protocol.String()
	} else if ip6Layer := packet.Layer(layers.LayerTypeIPv6); ip6Layer != nil {
		ip6 := ip6Layer.(*layers.IPv6)
		srcIP = ip6.SrcIP.String()
		dstIP = ip6.DstIP.String()
		protocol = ip6.NextHeader.String()
	} else {
		// Not an IP packet, ignore for this analyzer
		return AnalyzedPacket{}, false
	}

	// Extract Transport Layer (UDP or TCP)
	var srcPort, dstPort int
	isUDP := false

	if udpLayer := packet.Layer(layers.LayerTypeUDP); udpLayer != nil {
		udp := udpLayer.(*layers.UDP)
		srcPort = int(udp.SrcPort)
		dstPort = int(udp.DstPort)
		isUDP = true
	} else if tcpLayer := packet.Layer(layers.LayerTypeTCP); tcpLayer != nil {
		tcp := tcpLayer.(*layers.TCP)
		srcPort = int(tcp.SrcPort)
		dstPort = int(tcp.DstPort)
	} else {
		// No transport layer found, ignore
		return AnalyzedPacket{}, false
	}

	// Filter for "game-like" traffic:
	// Requirement: UDP traffic on a specific port range.
	if !isUDP {
		return AnalyzedPacket{}, false
	}

	// Check if either source or destination port falls within the filtered range.
	if (srcPort < a.filterPortMin || srcPort > a.filterPortMax) &&
		(dstPort < a.filterPortMin || dstPort > a.filterPortMax) {
		return AnalyzedPacket{}, false
	}

	return AnalyzedPacket{
		Timestamp:  p.Timestamp,
		SourceIP:   srcIP,
		DestIP:     dstIP,
		SourcePort: srcPort,
		DestPort:   dstPort,
		Protocol:   protocol,
		PayloadLen: len(p.Payload),
	}, true
}
