package service


// https://github.com/shurcooL/notifications
// Package fs implements notifications.Service using a virtual filesystem.
// Package githubapi implements notifications.Service using GitHub API clients.
// https://github.com/pachyderm/pachyderm

/*

import (
	// data
	"github.com/roscopecoltran/sniperkit-limo/model"
	// logs
	"github.com/sirupsen/logrus"
	prefixed "github.com/x-cray/logrus-prefixed-formatter"
	// notifications
	// "github.com/deckarep/gosx-notifier"
)

func gosxnotifierTest() {
    //At a minimum specifiy a message to display to end-user.
    note := gosxnotifier.NewNotification("Check your Apple Stock!")

    //Optionally, set a title
    note.Title = "It's money making time 💰"

    //Optionally, set a subtitle
    note.Subtitle = "My subtitle"

    //Optionally, set a sound from a predefined set.
    note.Sound = gosxnotifier.Basso

    //Optionally, set a group which ensures only one notification is ever shown replacing previous notification of same group id.
    note.Group = "com.unique.yourapp.identifier"

    //Optionally, set a sender (Notification will now use the Safari icon)
    note.Sender = "com.apple.Safari"

    //Optionally, specifiy a url or bundleid to open should the notification be
    //clicked.
    note.Link = "http://www.yahoo.com" //or BundleID like: com.apple.Terminal

    //Optionally, an app icon (10.9+ ONLY)
    note.AppIcon = "gopher.png"

    //Optionally, a content image (10.9+ ONLY)
    note.ContentImage = "gopher.png"

    //Then, push the notification
    err := note.Push()

    //If necessary, check error
    if err != nil {
        log.Println("Uh oh!")
    }
}
*/