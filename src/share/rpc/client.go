package rpc

import (
    "net"
    "strconv"
    "time"
    "share/lib/rpc2"
    "share/event"
)

// RPC Client struct
type Client struct {
    connected   bool
    conn        net.Conn
    rpcClient   *rpc2.Client

    IpAddress   string
    Port        int
}

// Initializes RPC Client
func (c *Client) Init() {
    log.Info("Attempting to connect to the Master Server...")

    c.connected = false
    go c.run()
}

/*
    Synchronous RPC Call method
    @param  method  function name
    @param  args    function call argument
    @param  reply   function return call result
    @return error, if any
 */
func (c *Client) Call(method string, args interface{}, reply interface{}) error {
    return c.rpcClient.Call(method, args, reply)
}


/*
    RPC Client run function, an infinite loop with 5 second delay.
    This will detect when connection will be closed and will attempt to reconnect
 */
func (c *Client) run() {
    for {
        if !c.connected {
            var conn, err = net.Dial("tcp", c.IpAddress + ":" + strconv.Itoa(c.Port))
            if err == nil {
                c.conn      = conn
                c.rpcClient = rpc2.NewClient(c.conn)
                c.connected = true
                go c.rpcClient.Run()
                event.Trigger(event.SyncConnectEvent, nil)
            }
        } else {
            var buffer= make([]byte, 0)

            _, err := c.conn.Read(buffer)
            if err != nil {
                c.conn.Close()
                c.connected = false
                event.Trigger(event.SyncDisconnectEvent, nil)
            }
        }

        time.Sleep(5 * time.Second)
    }
}