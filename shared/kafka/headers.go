package kafka

import "github.com/IBM/sarama"

func convertHeaders(headers []*sarama.RecordHeader) map[string]string {
	if len(headers) == 0 {
		return nil
	}

	m := make(map[string]string, len(headers))
	for _, h := range headers {
		if h == nil {
			continue
		}
		m[string(h.Key)] = string(h.Value)
	}
	return m
}
