package main

import (
	"context"
	"sort"
	"sync"
)

func newPostcodeFinder(ctx context.Context) *postcodeFinder {

	postcodeFinder := &postcodeFinder{
		&pipeline{
			wg:     sync.WaitGroup{},
			input:  make(chan delivery, 10),
			output: make(chan searchEngineResult, 2),
			ctx:    ctx,
		},
	}
	postcodeFinder.start()
	return postcodeFinder
}

type postcodeFinder struct {
	*pipeline
}

func (rc *postcodeFinder) sendData(d interface{}) {
	rc.pipeline.sendData(d)
}

func (rc *postcodeFinder) outputChannel() <-chan searchEngineResult {
	return rc.pipeline.output
}

func (rc *postcodeFinder) closeInputChannel() {
	rc.pipeline.closeInputChannel()
}

func (pf *postcodeFinder) start() {

	go func() {
		defer close(pf.output)

		deliveryByPostcod := make(map[string]int, MAX_POSSIBLE_POSTCODE)
		for {
			select {
			case <-pf.ctx.Done():
				return
			case order, ok := <-pf.input:
				if !ok {
					busiest := extractBusiest(&deliveryByPostcod)

					select {
					case <-pf.ctx.Done():
					case pf.output <- searchEngineResult{
						data: map[string]interface{}{
							"busiest_postcode": map[string]interface{}{
								"delivery_count": busiest.deliveryCount,
								"postcode":       busiest.postcode,
							},
						},
					}:
					}

					return
				}
				deliveryByPostcod[order.Postcode] += 1

			}
		}
	}()

}

func extractBusiest(deliveries *map[string]int) busiestPostcode {

	if len(*deliveries) == 0 {
		return busiestPostcode{}
	}

	deliveryList := make(postcodeList, 0, len(*deliveries))
	for k, v := range *deliveries {
		deliveryList = append(deliveryList, busiestPostcode{
			postcode:      k,
			deliveryCount: v,
		})
	}

	sort.Sort(deliveryList)
	return deliveryList[0]
}

type busiestPostcode struct {
	postcode      string
	deliveryCount int
}

type postcodeList []busiestPostcode

func (r postcodeList) Len() int {
	return len(r)
}

func (r postcodeList) Less(i, j int) bool {
	return r[i].deliveryCount > r[j].deliveryCount
}

func (r postcodeList) Swap(i, j int) {
	r[i], r[j] = r[j], r[i]
}
