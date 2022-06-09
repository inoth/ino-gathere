package host_server

type HostServer struct{}

func (hs *HostServer) ServeStart() error {
	select {}
}
