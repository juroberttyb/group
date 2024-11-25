package group

import (
	"context"
	"fmt"
	"sync"
	"testing"
	"time"
)

// test whether only one function would run while there are two threads spawned with the same key
// expected output
// run
// result <nil>
// result <nil>
func TestSameKeyOnlyOneRun(t *testing.T) {

	ctx := context.Background()
	wg := sync.WaitGroup{}
	wg.Add(2)
	group := NewGroup[string, any](Options{
		Timeout: 1 * time.Second,
	})
	go func() {
		defer wg.Done()
		fmt.Println(group.Do(ctx, "foo", someTask))
	}()
	go func() {
		defer wg.Done()
		fmt.Println(group.Do(ctx, "foo", someTask))
	}()
	wg.Wait()

	// meetupStore := new(store.MockMeetup)

	// layout := "2006-01-02 15:04:05"
	// createdAt, err := time.Parse(layout, "2024-01-22 04:06:20")
	// if err != nil {
	// 	t.Errorf("err: %s", err)
	// 	return
	// }
	// updatedAt, err := time.Parse(layout, "2024-01-22 04:06:20")
	// if err != nil {
	// 	t.Errorf("err: %s", err)
	// 	return
	// }
	// hostingAt, err := time.Parse(layout, "2025-01-22 05:38:39")
	// if err != nil {
	// 	t.Errorf("err: %s", err)
	// 	return
	// }
	// meetupStore.On("GetRedirect", mock.Anything).Return([]*models.Meetup{}, nil)
	// meetupStore.On("GetMeetups", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return([]*models.Meetup{
	// 	{
	// 		ID:               "7849583d-197c-48de-b48a-ce81cc26eca2",
	// 		CreatorID:        "b7c6fc25-0cc8-4b5b-a162-2d784fa9c0d9",
	// 		IsActive:         true,
	// 		IsDeleted:        false,
	// 		Status:           "attending",
	// 		Title:            "了不起的標題",
	// 		ParticipantCount: 0,
	// 		Tags:             []string{},
	// 		CreatedAt:        createdAt,
	// 		UpdatedAt:        updatedAt,
	// 		HostingAt:        hostingAt,
	// 	},
	// }, nil)
	// meetupStore.On("GetMyMeetupIDs", mock.Anything, mock.Anything, mock.Anything).Return([]string{}, nil)
	// userStore := new(store.MockUser)
	// userStore.On("Get", mock.Anything, mock.Anything).Return(&models.User{
	// 	Specialties: []string{"Prosthodontic"},
	// }, nil)
	// mediaStore := new(store.MockMedia)
	// storage := new(store.MockStorage)
	// notif := new(store.MockNotification)
	// mq := new(mq.MockMQ)
	// meetupSvc := NewMeetup(context.Background(), meetupStore, userStore, mediaStore, storage, notif, mq, mq, queue.NewPool(8))

	// meetups, _, err := meetupSvc.GetMeetups(context.Background(), "b3b646b7-7e37-4ed9-a4b3-11503b94763c", "", "", "", "", 10, models.Normal, false)
	// if err != nil {
	// 	t.Errorf("err: %s", err)
	// 	return
	// }
	// meetupStore.AssertExpectations(t)
	// userStore.AssertExpectations(t)

	// for _, m := range meetups {
	// 	require.Equal(t, m.Status, models.Normal, "status should be normal")
	// }
}

func TestConcurrentRunDiffKey(t *testing.T) {
	ctx := context.Background()
	wg := sync.WaitGroup{}
	wg.Add(2)
	group := NewGroup[string, any](Options{
		Timeout: 1 * time.Second,
	})
	go func() {
		defer wg.Done()
		fmt.Println(group.Do(ctx, "foo", someTask))
	}()

	go func() {
		defer wg.Done()
		fmt.Println(group.Do(ctx, "bar", someTask))
	}()
	wg.Wait()
}

func TestAggregateMeetupsAllStatus(t *testing.T) {
	// // Create a new context
	// ctx := context.Background()

	// // Create some meetups
	// meetups := []*models.Meetup{
	// 	{ID: "0", Tags: []string{"tag0", "tag1"}},
	// 	{ID: "1", Tags: []string{"specialty_0", "tag2"}, MeetupType: string(models.Internal)},
	// 	{ID: "2", Tags: []string{"specialty_0", "tag2"}},
	// 	{ID: "3", Tags: []string{"tag3"}},
	// 	{ID: "4", Tags: []string{"tag4"}},
	// }
	// uid := "b3b646b7-7e37-4ed9-a4b3-11503b94763c"

	// meetupStore := new(store.MockMeetup)
	// meetupStore.On("GetRedirect", mock.Anything).Return([]*models.Meetup{}, nil)
	// meetupStore.On("GetMyMeetupIDs", mock.Anything, uid, models.Interested).Return([]string{}, nil)
	// meetupStore.On("GetMyMeetupIDs", mock.Anything, uid, models.Attending).Return([]string{"0"}, nil)
	// meetupStore.On("GetMyMeetupIDs", mock.Anything, uid, models.Attended).Return([]string{"3"}, nil)
	// userStore := new(store.MockUser)
	// userStore.On("Get", mock.Anything, uid).Return(&models.User{
	// 	Specialties: []string{"specialty_0", "specialty_1"},
	// }, nil)
	// mediaStore := new(store.MockMedia)
	// storage := new(store.MockStorage)
	// notif := new(store.MockNotification)
	// mq := new(mq.MockMQ)
	// meetupSvc := NewMeetup(context.Background(), meetupStore, userStore, mediaStore, storage, notif, mq, mq, queue.NewPool(8))

	// // Call the aggregateMeetups function
	// meetups, err := meetupSvc.AggregateMeetups(ctx, meetups, uid, models.Normal, 0)
	// if err != nil {
	// 	t.Errorf("err: %s", err)
	// 	return
	// }

	// // Check the status of each meetup
	// require.Equal(t, models.Attending, meetups[0].Status)
	// require.Equal(t, models.Recommended, meetups[1].Status)
	// require.Equal(t, models.Normal, meetups[2].Status)
	// require.Equal(t, models.Attended, meetups[3].Status)
	// require.Equal(t, models.Normal, meetups[4].Status)
}
