package service

// ref. https://github.com/Jeffail/benthos

/*

frodo := tooth.New("frodo")
sam := tooth.New("sam")
gollum := tooth.New("gollum")

sam.Subscribe(frodo)
frodo.Subscribe(gollum)
sam.Subscribe(gollum)

frodo.Publish("I'm so tired")
gollum.Publish("My preciousssss")

msg1 := sam.FetchAll()      // [I'm so tired My preciousssss]
msg2 := frodo.Fetch(gollum) // My preciousssss

*/