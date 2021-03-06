/*
Copyright (c) 2015 VMware, Inc. All Rights Reserved.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package autostart

import (
	"errors"
	"flag"

	"github.com/juju/govmomi/govc/flags"
	"github.com/juju/govmomi/object"
	"github.com/juju/govmomi/vim25/methods"
	"github.com/juju/govmomi/vim25/mo"
	"github.com/juju/govmomi/vim25/types"
	"golang.org/x/net/context"
)

type AutostartFlag struct {
	*flags.ClientFlag
	*flags.DatacenterFlag
	*flags.HostSystemFlag
}

func (f *AutostartFlag) Register(fs *flag.FlagSet) {}

func (f *AutostartFlag) Process() error { return nil }

// VirtualMachines returns list of virtual machine objects based on the
// arguments specified on the command line. This helper is defined in
// flags.SearchFlag as well, but that pulls in other virtual machine flags that
// are not relevant here.
func (f *AutostartFlag) VirtualMachines(args []string) ([]*object.VirtualMachine, error) {
	if len(args) == 0 {
		return nil, errors.New("no argument")
	}

	finder, err := f.Finder()
	if err != nil {
		return nil, err
	}

	return finder.VirtualMachineList(context.TODO(), args...)
}

func (f *AutostartFlag) HostAutoStartManager() (*mo.HostAutoStartManager, error) {
	h, err := f.HostSystem()
	if err != nil {
		return nil, err
	}

	var mhs mo.HostSystem
	err = h.Properties(context.TODO(), h.Reference(), []string{"configManager.autoStartManager"}, &mhs)
	if err != nil {
		return nil, err
	}

	var mhas mo.HostAutoStartManager
	err = h.Properties(context.TODO(), *mhs.ConfigManager.AutoStartManager, nil, &mhas)
	if err != nil {
		return nil, err
	}

	return &mhas, nil
}

func (f *AutostartFlag) ReconfigureDefaults(template types.AutoStartDefaults) error {
	c, err := f.Client()
	if err != nil {
		return err
	}

	mhas, err := f.HostAutoStartManager()
	if err != nil {
		return err
	}

	req := types.ReconfigureAutostart{
		This: mhas.Reference(),
		Spec: types.HostAutoStartManagerConfig{
			Defaults: &template,
		},
	}

	_, err = methods.ReconfigureAutostart(context.TODO(), c, &req)
	if err != nil {
		return err
	}

	return nil
}

func (f *AutostartFlag) ReconfigureVMs(args []string, template types.AutoStartPowerInfo) error {
	c, err := f.Client()
	if err != nil {
		return err
	}

	mhas, err := f.HostAutoStartManager()
	if err != nil {
		return err
	}

	req := types.ReconfigureAutostart{
		This: mhas.Reference(),
		Spec: types.HostAutoStartManagerConfig{
			PowerInfo: make([]types.AutoStartPowerInfo, 0),
		},
	}

	vms, err := f.VirtualMachines(args)
	if err != nil {
		return err
	}

	for _, vm := range vms {
		pi := template
		pi.Key = vm.Reference()
		req.Spec.PowerInfo = append(req.Spec.PowerInfo, pi)
	}

	_, err = methods.ReconfigureAutostart(context.TODO(), c, &req)
	if err != nil {
		return err
	}

	return nil
}
