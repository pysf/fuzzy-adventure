package main

import (
	"context"
	"io"
	"reflect"
	"strings"
	"sync"
	"testing"
	"time"
)

type DeliveryReaderResult struct {
	input    downloadResult
	expected searchEnginInput
}

var DeliveryReaderResults = []DeliveryReaderResult{
	{
		input: downloadResult{
			f: io.NopCloser(strings.NewReader(`
			[
				{
					"postcode": "10118",
					"recipe": "Spanish One-Pan Chicken",
					"delivery": "Saturday 11AM - 1PM"
				}
			]
			`)),
		},
		expected: searchEnginInput{
			data: delivery{
				Postcode: "10118",
				Recipe:   "Spanish One-Pan Chicken",
				Delivery: "Saturday 11AM - 1PM",
			},
		},
	},
}

func TestDeliveryReaderStart(t *testing.T) {

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	deliveryReader := deliveryReader{
		ctx: ctx,
		wg:  sync.WaitGroup{},
	}

	downloadCh := make(chan downloadResult)
	resultCh := deliveryReader.start(downloadCh)

	for i, testExpect := range DeliveryReaderResults {

		go func() {
			defer close(downloadCh)
			downloadCh <- testExpect.input
		}()

		select {
		case resp, ok := <-resultCh:
			if !ok {
				break
			}

			if testExpect.expected.err != resp.err {
				t.Fatalf("case %d expeted %v got %v", i, testExpect.expected.err, resp.err)
			}

			if !reflect.DeepEqual(testExpect.expected.data, resp.data) {
				t.Fatalf("case %d expeted %v got %v", i, testExpect.expected.data, resp.data)
			}

		case <-time.After(1 * time.Second):
			t.Fatalf("%d ,timedout!", i)
		}

	}

}
