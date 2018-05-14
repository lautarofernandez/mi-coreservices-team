package jq

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"
)

const testcase = `{
    "id":"payment_v1_gateway-3c893ce4b9ebc57a3f67e94a0888806afe1c02b0",
    "internal_id":"payment_v1-3025573483-gateway",
    "type":"gateway",
    "site_id":"MLA",
    "user_id":265836623,
    "version":4,
    "schema_version":5,
    "schema_original_version":0,
    "date_created":"2017-10-02T00:42:08.000Z",
    "last_modified":"2018-04-19T13:45:09.209Z",
    "extra":{
        "has_notes":0,
        "actions":null,
        "shipping":{
            "status":"",
            "id":0
        },
        "payment":{
            "type":"regular_payment",
            "acquired":"amex",
            "payment_method_id":"amex",
            "merchant_number":"9909150666",
            "is_offline_payment":false
        },
        "mediation_allowed_to_create":false,
        "tags":""
    },
    "full_text":{
        "_all_data":"|3025573483|738225|5240|Ramirez|RAMA7478551|0342156156170|5240|fernanda@masturismo.tur.ar|Maria Fernanda|Lollapalooza 2018 - 18/03/2018 13:00 - 16/03/2018 13:00 - 17/03/2018 13:00|carlos osc friggeri",
        "amount":"5240",
        "card_name":"carlos osc friggeri",
        "email":"fernanda@masturismo.tur.ar",
        "external_reference":"738225",
        "first_name":"Maria Fernanda",
        "id":"3025573483",
        "last_name":"Ramirez",
        "nickname":"RAMA7478551",
        "phone_number":"0342156156170",
        "title":"Lollapalooza 2018 - 18/03/2018 13:00 - 16/03/2018 13:00 - 17/03/2018 13:00",
        "transaction_amount":"5240"
    },
    "row":{
        "amount":{
            "currency_code":"ARS",
            "value":5240
        },
        "counterpart":{
            "id":150275639,
            "nickname":"RAMA7478551",
            "email":"fernanda@masturismo.tur.ar",
            "first_name":"Maria Fernanda",
            "last_name":"Ramirez",
            "phones":{
                "main_phone":{
                    "area_code":" ",
                    "extension":"",
                    "number":"0342156156170"
                },
                "alternative_phone":{
                    "area_code":"",
                    "extension":"",
                    "number":""
                }
            },
            "extra":{
                "company":{
                    "brand_name":"",
                    "corporate_name":"",
                    "identification":""
                }
            }
        },
        "title":"Lollapalooza 2018 - 18/03/2018 13:00 - 16/03/2018 13:00 - 17/03/2018 13:00",
        "status":"approved",
        "status_detail":"partially_refunded",
        "thumbnail":""
    },
    "resources":{
        "main_resource":{
            "header":{
                "id":3025573483,
                "type":"payment_v1",
                "last_modified":"2018-04-17T19:30:33.000Z"
            },
            "object":{
                "acquirer":"amex",
                "acquirer_reconciliation":[
                    {
                        "authorization_code":null,
                        "batch_closing_date":null,
                        "batch_number":null,
                        "operation":"refund_capture",
                        "refund_id":373993807,
                        "terminal_number":null,
                        "transaction_number":null
                    }
                ],
                "additional_info":{
                    "items":[
                        {
                            "category_id":null,
                            "description":null,
                            "id":"1",
                            "picture_url":null,
                            "quantity":"1",
                            "title":"Costo del Envio",
                            "unit_price":"180"
                        },
                        {
                            "category_id":null,
                            "description":null,
                            "id":"1725",
                            "picture_url":null,
                            "quantity":"1",
                            "title":"3 Day Pass",
                            "unit_price":"5060"
                        }
                    ]
                },
                "api_version":"2",
                "application_id":8851253330337986,
                "authorization_code":"590000",
                "available_actions":[
                    "refund"
                ],
                "binary_mode":true,
                "call_for_authorize_id":null,
                "captured":true,
                "card":{
                    "cardholder":{
                        "identification":{
                            "number":"22715110",
                            "type":"DNI"
                        },
                        "name":"carlos osc friggeri"
                    },
                    "date_created":"2017-10-01T20:42:08.000-04:00",
                    "date_last_updated":"2017-10-01T20:42:08.000-04:00",
                    "expiration_month":4,
                    "expiration_year":2022,
                    "first_six_digits":"377792",
                    "id":null,
                    "last_four_digits":"2544"
                },
                "client_id":"8851253330337986",
                "collector":{
                    "email":"lolla2018@allaccess.com.ar",
                    "first_name":"DF ENTERTAINMENT",
                    "id":265836623,
                    "identification":{
                        "number":"30714942839",
                        "type":"Otro"
                    },
                    "last_name":"S.A.",
                    "phone":{
                        "area_code":null,
                        "extension":null,
                        "number":"111111111"
                    }
                },
                "collector_id":265836623,
                "collector_tags":[

                ],
                "connection":null,
                "counter_currency":null,
                "coupon_amount":0,
                "coupon_id":null,
                "currency_id":"ARS",
                "date_approved":"2017-10-01T20:42:11.000-04:00",
                "date_created":"2017-10-01T20:42:08.000-04:00",
                "date_last_updated":"2018-04-17T15:30:33.000-04:00",
                "date_of_expiration":null,
                "deduction_schema":null,
                "description":"Lollapalooza 2018 - 18/03/2018 13:00 - 16/03/2018 13:00 - 17/03/2018 13:00",
                "differential_pricing_id":null,
                "external_reference":"738225",
                "fee_details":[

                ],
                "financing_type":null,
                "id":3025573483,
                "installments":6,
                "internal_metadata":{

                },
                "issuer_id":"310",
                "live_mode":true,
                "marketplace":"NONE",
                "merchant_account_id":"afaf5b43718a1d08e594bf7855940ac8",
                "merchant_number":"9909150666",
                "merchant_services":{
                    "fraud_manual_review":false,
                    "fraud_scoring":true
                },
                "metadata":{

                },
                "money_release_date":"2017-10-03T20:42:11.000-04:00",
                "money_release_days":null,
                "money_release_schema":null,
                "notification_url":"https://api.boletius.com/paymentWs/mercadoPagoIPNV2",
                "operation_type":"regular_payment",
                "order":{

                },
                "payer":{
                    "email":"fernanda@masturismo.tur.ar",
                    "entity_type":null,
                    "first_name":"Maria Fernanda",
                    "id":null,
                    "identification":{
                        "number":null,
                        "type":null
                    },
                    "last_name":"Ramirez",
                    "phone":{
                        "area_code":" ",
                        "extension":null,
                        "number":"0342156156170"
                    },
                    "type":"guest"
                },
                "payer_id":150275639,
                "payer_tags":[

                ],
                "payment_method_id":"amex",
                "payment_type_id":"credit_card",
                "processing_mode":"gateway",
                "profile_id":"allaccess_gateway310amex_lolla2018",
                "refunds":[
                    {
                        "amount":1669.8,
                        "collector_movement_id":null,
                        "counter_currency":null,
                        "date_created":"2018-04-17T13:53:53.000-04:00",
                        "gtw_refund_id":6612432506,
                        "id":373993807,
                        "metadata":{

                        },
                        "payer_movement_id":null,
                        "payment_id":3025573483,
                        "source":{
                            "id":"265836623",
                            "name":"DF ENTERTAINMENT S.A.",
                            "type":"collector"
                        },
                        "status":null,
                        "unique_sequence_number":null
                    }
                ],
                "reserve_id":null,
                "risk_execution_id":14645401078,
                "shipping_amount":0,
                "site_id":"MLA",
                "sponsor_id":null,
                "statement_descriptor":"LOLLPALOOZA",
                "status":"approved",
                "status_detail":"partially_refunded",
                "transaction_amount":5240,
                "transaction_amount_refunded":1669.8,
                "transaction_details":{
                    "acquirer_reference":null,
                    "external_resource_url":null,
                    "financial_institution":null,
                    "installment_amount":873.33,
                    "net_received_amount":5240,
                    "overpaid_amount":0,
                    "payable_deferral_period":null,
                    "payment_method_reference_id":null,
                    "total_paid_amount":5240
                },
                "transaction_id":"3700164821_757a776d77786b797b7d"
            }
        },
        "other_resources":[
            {
                "header":{
                    "id":265836623,
                    "type":"user_collector",
                    "last_modified":"0001-01-01T00:00:00.000Z"
                },
                "object":{
                    "address":{
                        "address":"-",
                        "city":"-",
                        "state":"AR-C",
                        "zip_code":null
                    },
                    "alternative_phone":{
                        "area_code":"",
                        "extension":"",
                        "number":""
                    },
                    "bill_data":{
                        "accept_credit_note":"N"
                    },
                    "buyer_reputation":{
                        "canceled_transactions":0,
                        "tags":[

                        ],
                        "transactions":{
                            "canceled":{
                                "paid":null,
                                "total":null
                            },
                            "completed":null,
                            "not_yet_rated":{
                                "paid":null,
                                "total":null,
                                "units":null
                            },
                            "period":"historic",
                            "total":null,
                            "unrated":{
                                "paid":null,
                                "total":null
                            }
                        }
                    },
                    "company":{
                        "brand_name":"ALL ACCESS",
                        "city_tax_id":null,
                        "corporate_name":"DF ENTERTAINMENT S.A",
                        "identification":"30714942839",
                        "soft_descriptor":null,
                        "state_tax_id":null
                    },
                    "context":{
                        "device":"web-desktop",
                        "flow":"normal",
                        "source":"mercadopago"
                    },
                    "country_id":"AR",
                    "credit":{
                        "consumed":0,
                        "credit_level_id":"MLA1"
                    },
                    "email":"lolla2018@allaccess.com.ar",
                    "first_name":"DF ENTERTAINMENT",
                    "gender":"",
                    "id":265836623,
                    "identification":{
                        "number":"30714942839",
                        "type":"Otro"
                    },
                    "internal_tags":[
                        "gateway_mix"
                    ],
                    "last_name":"S.A.",
                    "logo":null,
                    "nickname":"DFENTERTAINMENTSAALLACCES",
                    "permalink":"http://perfil.mercadolibre.com.ar/DFENTERTAINMENTSAALLACCES",
                    "phone":{
                        "area_code":null,
                        "extension":"",
                        "number":"111111111",
                        "verified":false
                    },
                    "points":0,
                    "pwd_generation_status":"none",
                    "registration_date":"2017-07-26T08:53:32.000-04:00",
                    "seller_experience":"ADVANCED",
                    "seller_reputation":{
                        "level_id":null,
                        "metrics":{
                            "claims":{
                                "period":"60 months",
                                "rate":0
                            },
                            "delayed_handling_time":{
                                "period":"60 months",
                                "rate":0
                            },
                            "sales":{
                                "completed":0,
                                "period":"60 months"
                            }
                        },
                        "power_seller_status":null,
                        "transactions":{
                            "canceled":0,
                            "completed":0,
                            "period":"historic",
                            "ratings":{
                                "negative":0,
                                "neutral":0,
                                "positive":0
                            },
                            "total":0
                        }
                    },
                    "shipping_modes":[
                        "custom",
                        "not_specified"
                    ],
                    "site_id":"MLA",
                    "status":{
                        "billing":{
                            "allow":true,
                            "codes":[

                            ]
                        },
                        "buy":{
                            "allow":true,
                            "codes":[

                            ],
                            "immediate_payment":{
                                "reasons":[

                                ],
                                "required":false
                            }
                        },
                        "confirmed_email":false,
                        "immediate_payment":false,
                        "list":{
                            "allow":true,
                            "codes":[

                            ],
                            "immediate_payment":{
                                "reasons":[

                                ],
                                "required":false
                            }
                        },
                        "mercadoenvios":"not_accepted",
                        "mercadopago_account_type":"professional",
                        "mercadopago_tc_accepted":true,
                        "required_action":"",
                        "sell":{
                            "allow":true,
                            "codes":[

                            ],
                            "immediate_payment":{
                                "reasons":[

                                ],
                                "required":false
                            }
                        },
                        "shopping_cart":{
                            "buy":"allowed",
                            "sell":"allowed"
                        },
                        "site_status":"active",
                        "user_type":"simple_registration",
                        "verification_status":"POCF"
                    },
                    "tags":[
                        "normal",
                        "user_info_verified"
                    ],
                    "user_type":"normal"
                }
            },
            {
                "header":{
                    "id":150275639,
                    "type":"user_payer",
                    "last_modified":"0001-01-01T00:00:00.000Z"
                },
                "object":{
                    "address":{
                        "address":null,
                        "city":null,
                        "state":"AR-S",
                        "zip_code":null
                    },
                    "alternative_phone":{
                        "area_code":"",
                        "extension":"",
                        "number":""
                    },
                    "bill_data":{
                        "accept_credit_note":null
                    },
                    "buyer_reputation":{
                        "canceled_transactions":0,
                        "tags":[

                        ],
                        "transactions":{
                            "canceled":{
                                "paid":null,
                                "total":null
                            },
                            "completed":null,
                            "not_yet_rated":{
                                "paid":null,
                                "total":null,
                                "units":null
                            },
                            "period":"historic",
                            "total":null,
                            "unrated":{
                                "paid":null,
                                "total":null
                            }
                        }
                    },
                    "context":{
                        "device":"web-desktop",
                        "flow":"normal",
                        "source":"mercadolibre"
                    },
                    "country_id":"AR",
                    "credit":{
                        "consumed":0,
                        "credit_level_id":"MLA5"
                    },
                    "email":"fernanda@masturismo.tur.ar",
                    "first_name":"Maria Fernanda",
                    "gender":null,
                    "id":150275639,
                    "identification":{
                        "number":null,
                        "type":null
                    },
                    "last_name":"Ramirez",
                    "logo":null,
                    "nickname":"RAMA7478551",
                    "permalink":"http://perfil.mercadolibre.com.ar/RAMA7478551",
                    "phone":{
                        "area_code":" ",
                        "extension":"",
                        "number":"0342156156170",
                        "verified":false
                    },
                    "points":1,
                    "registration_date":"2017-04-11T18:51:34.000-04:00",
                    "secure_email":"fernan.s0wjb7@mail.mercadolibre.com",
                    "seller_experience":"NEWBIE",
                    "seller_reputation":{
                        "level_id":null,
                        "metrics":{
                            "claims":{
                                "period":"60 months",
                                "rate":0
                            },
                            "delayed_handling_time":{
                                "period":"60 months",
                                "rate":0
                            },
                            "sales":{
                                "completed":0,
                                "period":"60 months"
                            }
                        },
                        "power_seller_status":null,
                        "transactions":{
                            "canceled":0,
                            "completed":0,
                            "period":"historic",
                            "ratings":{
                                "negative":0,
                                "neutral":0,
                                "positive":0
                            },
                            "total":0
                        }
                    },
                    "shipping_modes":[
                        "custom",
                        "not_specified"
                    ],
                    "site_id":"MLA",
                    "status":{
                        "billing":{
                            "allow":false,
                            "codes":[
                                "address_pending",
                                "identification_pending",
                                "identification_min_length_not_satisfied",
                                "address_empty_city"
                            ]
                        },
                        "buy":{
                            "allow":true,
                            "codes":[

                            ],
                            "immediate_payment":{
                                "reasons":[

                                ],
                                "required":false
                            }
                        },
                        "confirmed_email":true,
                        "immediate_payment":false,
                        "list":{
                            "allow":false,
                            "codes":[
                                "address_pending",
                                "identification_pending",
                                "identification_min_length_not_satisfied",
                                "address_empty_city"
                            ],
                            "immediate_payment":{
                                "reasons":[

                                ],
                                "required":false
                            }
                        },
                        "mercadoenvios":"not_accepted",
                        "mercadopago_account_type":"personal",
                        "mercadopago_tc_accepted":true,
                        "required_action":"simple_registration",
                        "sell":{
                            "allow":true,
                            "codes":[

                            ],
                            "immediate_payment":{
                                "reasons":[

                                ],
                                "required":false
                            }
                        },
                        "shopping_cart":{
                            "buy":"not_allowed",
                            "sell":"allowed"
                        },
                        "site_status":"active",
                        "user_type":"eventual",
                        "verification_status":"NV"
                    },
                    "tags":[
                        "normal",
                        "messages_as_buyer"
                    ],
                    "user_type":"normal"
                }
            }
        ]
    }
}`

func TestParseRules(t *testing.T) {
	rules := []string{
		"id",
		"internal_id",
		"type",
		"site_id",
		"user_id",
		"version",
		"schema_version",
		"schema_original_version",
		"date_created",
		"last_modified",

		"extra.actions",
		"extra.payment",

		"row.title",

		"resources.main_resource.header",
		"resources.other_resources[].header.id",
		"resources.other_resources[].object.address",
		"resources.other_resources[].not.exists",
	}

	ruleset := ParseRules(rules)

	// Manually test 2 levels of rule parsing results
	require.Len(t, ruleset, 13)
	for _, r := range ruleset {
		switch r.key {
		case "extra":
			require.Len(t, *r.child, 2)
			require.False(t, false, r.arrayChild)
		case "row":
			require.Len(t, *r.child, 1)
			require.False(t, false, r.arrayChild)
		case "resources":
			require.False(t, false, r.arrayChild)
			for _, r := range *r.child {
				switch r.key {
				case "main_resource":
					require.Len(t, *r.child, 1)
					require.False(t, false, r.arrayChild)
				case "other_resource":
					require.Len(t, *r.child, 3)
					require.False(t, true, r.arrayChild)
				}
			}
		}
	}
}

func TestRulesetFilter(t *testing.T) {
	tt := []struct {
		Name     string
		Input    string
		Expected string
		Rules    []string
	}{
		{
			Name:     "Simple JSON",
			Input:    `{"id": 132456, "data": {"amount": 12.12, "quantity": 5}, "type": "payment"}`,
			Expected: `{"id": 132456, "data": {"amount": 12.12}}`,
			Rules:    []string{"id", "data.amount"},
		},
		{
			Name:     "Simple JSON With Array",
			Input:    `{"id": 132456, "data": [{"amount": 12.12, "quantity": 5}, {"amount": 9.72, "quantity": 1}], "type": "payment"}`,
			Expected: `{"id": 132456, "data": [{"amount": 12.12}, {"amount":9.72}]}`,
			Rules:    []string{"id", "data[].amount"},
		},
		{
			Name:     "Complex JSON",
			Input:    testcase,
			Expected: `{"date_created":"2017-10-02T00:42:08.000Z","extra":{"actions":null,"payment":{"acquired":"amex","is_offline_payment":false,"merchant_number":"9909150666","payment_method_id":"amex","type":"regular_payment"}},"id":"payment_v1_gateway-3c893ce4b9ebc57a3f67e94a0888806afe1c02b0","internal_id":"payment_v1-3025573483-gateway","last_modified":"2018-04-19T13:45:09.209Z","resources":{"main_resource":{"header":{"id":3025573483,"last_modified":"2018-04-17T19:30:33.000Z","type":"payment_v1"}},"other_resources":[{"header":{"id":265836623},"object":{"address":{"state":"AR-C"}}},{"header":{"id":150275639},"object":{"address":{"state":"AR-S"}}}]},"row":{"title":"Lollapalooza 2018 - 18/03/2018 13:00 - 16/03/2018 13:00 - 17/03/2018 13:00"},"schema_original_version":0,"schema_version":5,"site_id":"MLA","type":"gateway","user_id":265836623,"version":4}`,
			Rules: []string{
				"id",
				"internal_id",
				"type",
				"site_id",
				"user_id",
				"version",
				"schema_version",
				"schema_original_version",
				"date_created",
				"last_modified",

				"extra.actions",
				"extra.payment",

				"row.title",

				"resources.main_resource.header",
				"resources.other_resources[].header.id",
				"resources.other_resources[].object.address.state",
				"resources.other_resources[].not.exists",
			},
		},
	}

	for _, tc := range tt {
		t.Run(tc.Name, func(t *testing.T) {
			var obj map[string]interface{}
			require.NoError(t, json.Unmarshal([]byte(tc.Input), &obj))

			ruleset := ParseRules(tc.Rules)

			filtered := ruleset.Filter(obj)
			out, err := json.Marshal(filtered)
			require.NoError(t, err)

			require.JSONEq(t, tc.Expected, string(out))
		})
	}
}

func BenchmarkRulesetFilter(b *testing.B) {
	var obj map[string]interface{}
	json.Unmarshal([]byte(testcase), &obj)

	rules := []string{
		"id",
		"internal_id",
		"type",
		"site_id",
		"user_id",
		"version",
		"schema_version",
		"schema_original_version",
		"date_created",
		"last_modified",

		"extra.actions",
		"extra.payment",

		"row.title",

		"resources.main_resource.header",
		"resources.other_resources[].header",
	}

	ruleset := ParseRules(rules)

	b.ResetTimer()

	for n := 0; n < b.N; n++ {
		ruleset.Filter(obj)
	}
}
