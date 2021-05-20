package quickbooks

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

// Purchase represents a QuickBooks Purchase object.
type Purchase struct {
	ID                   string        `json:"Id,omitempty"`
	DocNumber            string        `json:",omitempty"`
	TotalAmt             json.Number   `json:",omitempty"`
	TxnDate              Date          `json:",omitempty"`
	PaymentType          string        `json:",omitempty"`
	PaymentMethodRef     ReferenceType `json:",omitempty"`
	EntityRef            ReferenceType `json:",omitempty"`
	AccountRef           ReferenceType `json:",omitempty"`
	PrivateNote          string        `json:",omitempty"`
	GlobalTaxCalculation string        `json:",omitempty"`
	Line                 []Line        `json:",omitempty"`
}

// CreatePurchase creates the given Purchase on the QuickBooks server, returning
// the resulting Purchase object.
func (c *Client) CreatePurchase(pur *Purchase) (*Purchase, error) {
	var u, err = url.Parse(string(c.Endpoint))
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	u.Path = "/v3/company/" + c.RealmID + "/purchase"
	var v = url.Values{}
	v.Add("minorversion", minorVersion)
	u.RawQuery = v.Encode()
	var j []byte
	j, err = json.Marshal(pur)
	if err != nil {
		return nil, err
	}
	var req *http.Request
	req, err = http.NewRequest("POST", u.String(), bytes.NewBuffer(j))
	if err != nil {
		return nil, err
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")
	var res *http.Response
	res, err = c.Client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, parseFailure(res)
	}

	var r struct {
		Purchase Purchase
		Time     Date
	}
	err = json.NewDecoder(res.Body).Decode(&r)
	return &r.Purchase, err
}

/*

{
  "SalesReceipt": {
    "DocNumber": "1003",
    "SyncToken": "0",
    "domain": "QBO",
    "Balance": 0,
    "PaymentMethodRef": {
      "name": "Check",
      "value": "2"
    },
    "BillAddr": {
      "Lat": "INVALID",
      "Long": "INVALID",
      "Id": "49",
      "Line1": "Dylan Sollfrank"
    },
    "DepositToAccountRef": {
      "name": "Checking",
      "value": "35"
    },
    "TxnDate": "2014-09-14",
    "TotalAmt": 337.5,
    "CustomerRef": {
      "name": "Dylan Sollfrank",
      "value": "6"
    },
    "CustomerMemo": {
      "value": "Thank you for your business and have a great day!"
    },
    "PrintStatus": "NotSet",
    "PaymentRefNum": "10264",
    "EmailStatus": "NotSet",
    "sparse": false,
    "Line": [
      {
        "Description": "Custom Design",
        "DetailType": "SalesItemLineDetail",
        "SalesItemLineDetail": {
          "TaxCodeRef": {
            "value": "NON"
          },
          "Qty": 4.5,
          "UnitPrice": 75,
          "ItemRef": {
            "name": "Design",
            "value": "4"
          }
        },
        "LineNum": 1,
        "Amount": 337.5,
        "Id": "1"
      },
      {
        "DetailType": "SubTotalLineDetail",
        "Amount": 337.5,
        "SubTotalLineDetail": {}
      }
    ],
    "ApplyTaxAfterDiscount": false,
    "CustomField": [
      {
        "DefinitionId": "1",
        "Type": "StringType",
        "Name": "Crew #"
      }
    ],
    "Id": "11",
    "TxnTaxDetail": {
      "TotalTax": 0
    },
    "MetaData": {
      "CreateTime": "2014-09-16T14:59:48-07:00",
      "LastUpdatedTime": "2014-09-16T14:59:48-07:00"
    }
  },
  "time": "2015-07-29T09:29:56.229-07:00"
}
*/
