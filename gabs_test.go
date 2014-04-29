package gabs

import (
	"testing"
)

func TestBasic(t *testing.T) {
	sample := []byte(`{"test":{"value":10},"test2":20}`)

	val, err := ParseJson(sample)
	if err != nil {
		t.Errorf("Failed to parse: %v", err)
		return
	}

	if result, ok := val.Search([]string{"test", "value"}...).Data().(float64); ok {
		if result != 10 {
			t.Errorf("Wrong value of result: %v", result)
		}
	} else {
		t.Errorf("Didn't find test.value")
	}

	if _, ok := val.Search("test2", "value").Data().(string); ok {
		t.Errorf("Somehow found a field that shouldn't exist")
	}

	if result, ok := val.Search("test2").Data().(float64); ok {
		if result != 20 {
			t.Errorf("Wrong value of result: %v", result)
		}
	} else {
		t.Errorf("Didn't find test2")
	}
}

func TestModify(t *testing.T) {
	sample := []byte(`{"test":{"value":10},"test2":20}`)

	val, err := ParseJson(sample)
	if err != nil {
		t.Errorf("Failed to parse: %v", err)
		return
	}

	if err := val.Set(45.0, "test", "value"); err != nil {
		t.Errorf("Failed to set field")
	}

	if result, ok := val.Search([]string{"test", "value"}...).Data().(float64); ok {
		if result != 45 {
			t.Errorf("Wrong value of result: %v", result)
		}
	} else {
		t.Errorf("Didn't find test.value")
	}
}

func TestArrays(t *testing.T) {
	json1, _ := ParseJson([]byte(`{
		"languages":{
			"english":{
				"places":0
			},
			"french": {
				"places": [
					"france",
					"belgium"
				]
			}
		}
	}`))

	json2, _ := ParseJson([]byte(`{
		"places":[
			"great_britain",
			"united_states_of_america",
			"the_world"
		]
	}`))

	if english_places := json2.Search("places").Data(); english_places != nil {
		json1.Set(english_places, "languages", "english", "places")
	} else {
		t.Errorf("Didn't find places in json2")
	}

	if english_places := json1.Search("languages", "english", "places").Data(); english_places != nil {

		english_array, ok := english_places.([]interface{})
		if !ok {
			t.Errorf("places in json1 (%v) was not an array", english_places)
		}

		if len(english_array) != 3 {
			t.Errorf("wrong length of array: %v", len(english_array))
		}

	} else {
		t.Errorf("Didn't find places in json1")
	}
}

func TestLargeSample(t *testing.T) {
	sample := []byte(`{
		"test":{
			"innerTest":{
				"value":10,
				"value2":22,
				"value3":{
					"moreValue":45
				}
			}
		},
		"test2":20
	}`)

	val, err := ParseJson(sample)
	if err != nil {
		t.Errorf("Failed to parse: %v", err)
		return
	}

	if result, ok := val.Search("test", "innerTest", "value3", "moreValue").Data().(float64); ok {
		if result != 45 {
			t.Errorf("Wrong value of result: %v", result)
		}
	} else {
		t.Errorf("Didn't find value")
	}
}