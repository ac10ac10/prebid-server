package schain

import (
	"encoding/json"
	"testing"

	"github.com/mxmCherry/openrtb/v15/openrtb2"
	"github.com/prebid/prebid-server/openrtb_ext"
	"github.com/stretchr/testify/assert"
)

func TestSChainWriter(t *testing.T) {

	const seller1SChain string = `"schain":{"complete":1,"nodes":[{"asi":"directseller1.com","sid":"00001","rid":"BidRequest1","hp":1}],"ver":"1.0"}`
	const seller2SChain string = `"schain":{"complete":2,"nodes":[{"asi":"directseller2.com","sid":"00002","rid":"BidRequest2","hp":2}],"ver":"2.0"}`
	const seller3SChain string = `"schain":{"complete":3,"nodes":[{"asi":"directseller3.com","sid":"00003","rid":"BidRequest3","hp":3}],"ver":"3.0"}`
	const sellerWildCardSChain string = `"schain":{"complete":1,"nodes":[{"asi":"wildcard1.com","sid":"wildcard1","rid":"WildcardReq1","hp":1}],"ver":"1.0"}`

	tests := []struct {
		description string
		giveRequest openrtb2.BidRequest
		giveBidder  string
		wantRequest openrtb2.BidRequest
		wantError   bool
	}{
		{
			description: "nil source and nil ext.prebid.schains",
			giveRequest: openrtb2.BidRequest{
				Ext:    nil,
				Source: nil,
			},
			giveBidder: "appnexus",
			wantRequest: openrtb2.BidRequest{
				Ext:    nil,
				Source: nil,
			},
		},
		{
			description: "Use source schain -- no bidder schain or wildcard schain in nil ext.prebid.schains",
			giveRequest: openrtb2.BidRequest{
				Ext: json.RawMessage(`{}`),
				Source: &openrtb2.Source{
					Ext: json.RawMessage(`{` + seller2SChain + `}`),
				},
			},
			giveBidder: "appnexus",
			wantRequest: openrtb2.BidRequest{
				Ext: json.RawMessage(`{}`),
				Source: &openrtb2.Source{
					Ext: json.RawMessage(`{` + seller2SChain + `}`),
				},
			},
		},
		{
			description: "Use source schain -- no bidder schain or wildcard schain in not nil ext.prebid.schains",
			giveRequest: openrtb2.BidRequest{
				Ext: json.RawMessage(`{"prebid":{"schains":[{"bidders":["appnexus"],` + seller1SChain + `}]}}`),
				Source: &openrtb2.Source{
					Ext: json.RawMessage(`{` + seller2SChain + `}`),
				},
			},
			giveBidder: "rubicon",
			wantRequest: openrtb2.BidRequest{
				Ext: json.RawMessage(`{"prebid":{"schains":[{"bidders":["appnexus"],` + seller1SChain + `}]}}`),
				Source: &openrtb2.Source{
					Ext: json.RawMessage(`{` + seller2SChain + `}`),
				},
			},
		},
		{
			description: "Use schain for bidder in ext.prebid.schains; ensure other ext.source field values are retained.",
			giveRequest: openrtb2.BidRequest{
				Ext: json.RawMessage(`{"prebid":{"schains":[{"bidders":["appnexus"],` + seller1SChain + `}]}}`),
				Source: &openrtb2.Source{
					FD:     1,
					TID:    "tid data",
					PChain: "pchain data",
					Ext:    json.RawMessage(`{` + seller2SChain + `}`),
				},
			},
			giveBidder: "appnexus",
			wantRequest: openrtb2.BidRequest{
				Ext: json.RawMessage(`{"prebid":{"schains":[{"bidders":["appnexus"],` + seller1SChain + `}]}}`),
				Source: &openrtb2.Source{
					FD:     1,
					TID:    "tid data",
					PChain: "pchain data",
					Ext:    json.RawMessage(`{` + seller1SChain + `}`),
				},
			},
		},
		{
			description: "Use schain for bidder in ext.prebid.schains, nil req.source ",
			giveRequest: openrtb2.BidRequest{
				Ext:    json.RawMessage(`{"prebid":{"schains":[{"bidders":["appnexus"],` + seller1SChain + `}]}}`),
				Source: nil,
			},
			giveBidder: "appnexus",
			wantRequest: openrtb2.BidRequest{
				Ext: json.RawMessage(`{"prebid":{"schains":[{"bidders":["appnexus"],` + seller1SChain + `}]}}`),
				Source: &openrtb2.Source{
					Ext: json.RawMessage(`{` + seller1SChain + `}`),
				},
			},
		},
		{
			description: "Use wildcard schain in ext.prebid.schains.",
			giveRequest: openrtb2.BidRequest{
				Ext: json.RawMessage(`{"prebid":{"schains":[{"bidders":["*"],` + sellerWildCardSChain + `}]}}`),
				Source: &openrtb2.Source{
					Ext: nil,
				},
			},
			giveBidder: "appnexus",
			wantRequest: openrtb2.BidRequest{
				Ext: json.RawMessage(`{"prebid":{"schains":[{"bidders":["*"],` + sellerWildCardSChain + `}]}}`),
				Source: &openrtb2.Source{
					Ext: json.RawMessage(`{` + sellerWildCardSChain + `}`),
				},
			},
		},
		{
			description: "Use schain for bidder in ext.prebid.schains instead of wildcard.",
			giveRequest: openrtb2.BidRequest{
				Ext: json.RawMessage(`{"prebid":{"schains":[{"bidders":["appnexus"],` + seller1SChain + `},{"bidders":["*"],` + sellerWildCardSChain + `}]}}`),
				Source: &openrtb2.Source{
					Ext: nil,
				},
			},
			giveBidder: "appnexus",
			wantRequest: openrtb2.BidRequest{
				Ext: json.RawMessage(`{"prebid":{"schains":[{"bidders":["appnexus"],` + seller1SChain + `},{"bidders":["*"],` + sellerWildCardSChain + `}]}}`),
				Source: &openrtb2.Source{
					Ext: json.RawMessage(`{` + seller1SChain + `}`),
				},
			},
		},
		{
			description: "Use source schain -- multiple (two) bidder schains in ext.prebid.schains.",
			giveRequest: openrtb2.BidRequest{
				Ext: json.RawMessage(`{"prebid":{"schains":[{"bidders":["appnexus"],` + seller1SChain + `},{"bidders":["appnexus"],` + seller2SChain + `}]}}`),
				Source: &openrtb2.Source{
					Ext: json.RawMessage(`{` + seller3SChain + `}`),
				},
			},
			giveBidder: "appnexus",
			wantRequest: openrtb2.BidRequest{
				Ext: json.RawMessage(`{"prebid":{"schains":[{"bidders":["appnexus"],` + seller1SChain + `},{"bidders":["appnexus"],` + seller2SChain + `}]}}`),
				Source: &openrtb2.Source{
					Ext: json.RawMessage(`{` + seller3SChain + `}`),
				},
			},
			wantError: true,
		},
	}

	for _, tt := range tests {
		// unmarshal ext to get schains object needed to initialize writer
		var reqExt *openrtb_ext.ExtRequest
		if tt.giveRequest.Ext != nil {
			reqExt = &openrtb_ext.ExtRequest{}
			err := json.Unmarshal(tt.giveRequest.Ext, reqExt)
			if err != nil {
				t.Error("Unable to unmarshal request.ext")
			}
		}

		writer, err := NewSChainWriter(reqExt)

		if tt.wantError {
			assert.NotNil(t, err)
			assert.Nil(t, writer)
		} else {
			assert.Nil(t, err)
			assert.NotNil(t, writer)

			writer.Write(&tt.giveRequest, tt.giveBidder)

			assert.Equal(t, tt.wantRequest, tt.giveRequest, tt.description)
		}
	}
}
