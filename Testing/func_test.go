package main

import (
	"testing"
	"time"
	"fmt"
)

type testData struct {
	playerID        string
	scheduleDays    ScheduleDays
	timeZoneOffset  int
	currentTime     time.Time
	expectedPass    bool
	expectedError   error
}

func TestUserTimeToDeliver(t *testing.T) {
	testCases := []testData{
		{
			playerID:        "player1",
			scheduleDays:    ScheduleDays{"Monday", "Tuesday", "Wednesday"},
			timeZoneOffset:  -6,
			currentTime:     time.Date(2022, time.January, 3, 16, 0, 0, 0, time.UTC), //Monday 4:00pm
			expectedPass:    true,
			expectedError:   nil,
		},
		{
			playerID:        "player2",
			scheduleDays:    ScheduleDays{"Monday", "Tuesday", "Wednesday"},
			timeZoneOffset:  -6,
			currentTime:     time.Date(2022, time.January, 8, 16, 0, 0, 0, time.UTC), //Saturday 4:00pm
			expectedPass:    false,
			expectedError:   nil,
		},
		{
			playerID:        "",
			scheduleDays:    ScheduleDays{"Monday", "Tuesday", "Wednesday"},
			timeZoneOffset:  -6,
			currentTime:     time.Date(2022, time.January, 3, 16, 0, 0, 0, time.UTC), //Monday 4:00pm
			expectedPass:    false,
			expectedError:   fmt.Errorf("invalid player ID"),
		},
		{
			playerID:        "player4",
			scheduleDays:    ScheduleDays{},
			timeZoneOffset:  -6,
			currentTime:     time.Date(2022, time.January, 3, 16, 0, 0, 0, time.UTC), //Monday 4:00pm
			expectedPass:    true,
			expectedError:   nil,
		},
	}

	
}

		





							// TESTING 8
// func TestUserCurrentTime(t *testing.T) {
// 	testCases := []struct {
// 		name         string
// 		timeZone     int
// 		expectedTime string
// 	}{
// 		{
// 			name:         "New York",
// 			timeZone:     -5,
// 			expectedTime: time.Now().UTC().Add(time.Hour * -5).Format("03pm"),
// 		},
// 		{
// 			name:         "London",
// 			timeZone:     0,
// 			expectedTime: time.Now().UTC().Add(time.Hour * 0).Format("03pm"),
// 		},
// 		{
// 			name:         "Bangkok",
// 			timeZone:     7,
// 			expectedTime: time.Now().UTC().Add(time.Hour * 7).Format("03pm"),
// 		},
// 		{
// 			name:         "Dhaka",
// 			timeZone:     7,
// 			expectedTime: time.Now().UTC().Add(time.Hour * 7).Format("03pm"),
// 		},
// 	}
// 	for _, tc := range testCases {
// 		t.Run(tc.name, func(t *testing.T) {
// 			currentTime := userCurrentTime(tc.timeZone)
// 			assert.Equal(t, tc.expectedTime, currentTime)
// 		})
// 	}
// }



							// TEST 7

// type testCase struct {
// 	name string
// 	args args
// 	want int
// }

// type args struct {
// 	timeZoneOffset int
// }

// func TestUserCurrentHour(t *testing.T) {
// 	testCases := []testCase{
// 		{
// 			name: "returns current hour in same day when time zone offset is negative",
// 			args: args{
// 				timeZoneOffset: -5,
// 			},
// 			want: time.Now().UTC().Add(time.Hour * -5).Hour(),
// 		},
// 		{
// 			name: "returns current hour in next day when time zone offset is negative",
// 			args: args{
// 				timeZoneOffset: -13,
// 			},
// 			want: time.Now().UTC().Add(time.Hour * -13).Hour(),
// 		},
// 		{
// 			name: "returns current hour in same day when time zone offset is positive",
// 			args: args{
// 				timeZoneOffset: 5,
// 			},
// 			want: time.Now().UTC().Add(time.Hour * 5).Hour(),
// 		},
// 		// Add more test cases here
// 	}

// 	for _, tc := range testCases {
// 		t.Run(tc.name, func(t *testing.T) {
// 			localHour := userCurrentHour(tc.args.timeZoneOffset)
// 			assert.Equal(t, tc.want, localHour)
// 		})
// 	}
// }
										// TEST 6

// type testCase struct {
// 	name string
// 	args args
// 	want int
// }

// type args struct {
// 	hr             int
// 	timeZoneOffset int
// }

// func TestUserLocalHour(t *testing.T) {
// 	testCases := []testCase{
// 		// PASS
// 		{
// 			name: "returns local hour in same day when time zone offset is negative",
// 			args: args{
// 				hr:             12,
// 				timeZoneOffset: -5,
// 			},
// 			want: 7,
// 		},
// 		// FAIL
// 		{
// 			name: "returns local hour in next day when time zone offset is negative",
// 			args: args{
// 				hr:             12,
// 				timeZoneOffset: -13,
// 			},
// 			want: 11,
// 		},
// 		// PASS
// 		{
// 			name: "returns local hour in same day when time zone offset is positive",
// 			args: args{
// 				hr:             12,
// 				timeZoneOffset: 5,
// 			},
// 			want: 17,
// 		},
// 	}

// 	for _, tc := range testCases {
// 		t.Run(tc.name, func(t *testing.T) {
// 			localHour := userLocalHour(tc.args.hr, tc.args.timeZoneOffset)
// 			assert.Equal(t, tc.want, localHour)
// 		})
// 	}
// }

								// TESTING 5

// type testCase struct {
// 	name string
// 	args args
// }

// type args struct {
// 	index       int
// 	expectedZon string
// }

// func TestGetUserTime(t *testing.T) {
// 	testCases := []testCase{
// 		{
// 			name: "returns empty string for negative index",
// 			args: args{
// 				index:       -1,
// 				expectedZon: "",
// 			},
// 		},
// 		{
// 			name: "returns timezone for valid index",
// 			args: args{
// 				index:       0,
// 				expectedZon: "Prime",
// 			},
// 		},
// 		{
// 			name: "returns timezone for valid index",
// 			args: args{
// 				index:       5,
// 				expectedZon: "User",
// 			},
// 		},
// 		// Add more test cases here
// 	}

// 	for _, tc := range testCases {
// 		t.Run(tc.name, func(t *testing.T) {
// 			zone := getUserTime(tc.args.index)
// 			assert.Equal(t, tc.args.expectedZon, zone)
// 		})
// 	}
// }



						//TESTING 4

// type testCase struct {
// 	name string
// 	args args
// }

// type args struct {
// 	index       int
// 	expectedZon string
// }

// func TestGetPrimeTime(t *testing.T) {
// 	testCases := []testCase{
// 		{
// 			name: "returns empty string for negative index",
// 			args: args{
// 				index:       -1,
// 				expectedZon: "",
// 			},
// 		},
// 		{
// 			name: "returns lowercase time of day for valid index",
// 			args: args{
// 				index:       0,
// 				expectedZon: "best",
// 			},
// 		},
// 		{
// 			name: "returns lowercase time of day for valid index",
// 			args: args{
// 				index:       5,
// 				expectedZon: "12am",
// 			},
// 		},

// 	}

// 	for _, tc := range testCases {
// 		t.Run(tc.name, func(t *testing.T) {
// 			zone := getPrimeTime(tc.args.index)
// 			assert.Equal(t, tc.args.expectedZon, zone)
// 		})
// 	}
// }

					// TESTING 3
// type testCase struct {
// 	name string
// 	args args
// }

// type args struct {
// 	index       int
// 	expectedZon string
// }

// func TestGetScheduleTimeOfDay(t *testing.T) {
// 	testCases := []testCase{
// 		{
// 			name: "returns empty string for out of range index",
// 			args: args{
// 				index:       -1,
// 				expectedZon: "",
// 			},
// 		},
// 		{
// 			name: "returns time of day for valid index",
// 			args: args{
// 				index:       0,
// 				expectedZon: "Best",
// 			},
// 		},
// 		{
// 			name: "returns time of day for valid index",
// 			args: args{
// 				index:       5,
// 				expectedZon: "12am",
// 			},
// 		},
// 	}

// 	for _, tc := range testCases {
// 		t.Run(tc.name, func(t *testing.T) {
// 			zone := getScheduleTimeOfDay(tc.args.index)
// 			assert.Equal(t, tc.args.expectedZon, zone)
// 		})
// 	}
// }


							// TESTING 2
// type args struct {
// 	index       int
// 	expected string
// }

// type testCase struct {
// 	name string
// 	args args
// }

// func TestGetScheduleTimeZone(t *testing.T) {
// 	testCases := []testCase{
// 		// PASS
// 		{
// 			name: "returns empty string for out of range index",
// 			args: args{
// 				index:       -1,
// 				expected: "",
// 			},
// 		},
// 		// PASS
// 		{
// 			name: "returns timezone for valid index",
// 			args: args{
// 				index:       5,
// 				expected: "User",
// 			},
// 		},
// 		// FAILL
// 		{
// 			name: "returns timezone for valid index",
// 			args: args{
// 				index:       0,
// 				expected: "Prime Time Best", 	// prime
// 			},
// 		},
// 	}

// 	for _, tc := range testCases {
// 		t.Run(tc.name, func(t *testing.T) {
// 			zone := getScheduleTimeZone(tc.args.index)
// 			assert.Equal(t, tc.args.expected, zone)
// 		})
// 	}
// }



						// TESTING 1

// func TestScheduleTimes(t *testing.T) {
// 	expectedTimes := []string{
// 		"Prime Time Best", 
// 		"Prime Time Morning", "Prime Time Afternoon", "Prime Time Evening", "Prime Time Swing", "User 12am", "User 01am", "User 02am", "User 03am", "User 04am", "User 05am", "User 06am", "User 07am", "User 08am", "User 09am", "User 10am", "User 11am", "User 12pm", "User 01pm", "User 02pm", "User 03pm", "User 04pm", "User 05pm", "User 06pm", "User 07pm", "User 08pm", "User 09pm", "User 10pm", "User 11pm", "PT 12am", "PT 01am", "PT 02am", "PT 03am", "PT 04am", "PT 05am", "PT 06am", "PT 07am", "PT 08am", "PT 09am", "PT 10am", "PT 11am", "PT 12pm", "PT 01pm", "PT 02pm", "PT 03pm", "PT 04pm", "PT 05pm", "PT 06pm", "PT 07pm", "PT 08pm", "PT 09pm", "PT 10pm", "PT 11pm", "GMT 12am", "GMT 01am", "GMT 02am", "GMT 03am", "GMT 04am", "GMT 05am", "GMT 06am", "GMT 07am", "GMT 08am", "GMT 09am", "GMT 10am", "GMT 11am", "GMT 12pm", "GMT 01pm", "GMT 02pm", 
// 		"GMT 03pm", "GMT 04pm", "GMT 05pm", "GMT 06pm", "GMT 07pm", "GMT 08pm", "GMT 09pm", "GMT 10pm", "GMT 11pm",
// 	}
// 	times := scheduleTimes()
// 	assert.Equal(t, expectedTimes, times)
// }


