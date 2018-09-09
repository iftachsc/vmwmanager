package main

import (
	"context"

	"github.com/iftachsc/vmware"
	"github.com/vmware/govmomi"
)

type VmwareVcenter struct {
	Host     string
	User     string
	Password string
}

func (vc VmwareVcenter) GetHost() string {
	return vc.Host
}

func (vc VmwareVcenter) GetUser() string {
	return vc.User
}
func (vc VmwareVcenter) GetPassword() string {
	return vc.Password
}

func (vc VmwareVcenter) InitializeClient(ctx context.Context) (*govmomi.Client, error) {
	println("INITIALIZING VM CLIENT")
	return vmware.NewClient(ctx, vc.Host, vc.User, vc.Password)
	//return vmware.NewClientFromEnv(ctx)
}
