package v20170131

import "testing"

func TestMerge(t *testing.T) {
	newCS := &ContainerService{
		Properties: &Properties{
			ServicePrincipalProfile: &ServicePrincipalProfile{
				ClientID: "fakeID",
				Secret:   "",
			},
			WindowsProfile: &WindowsProfile{
				AdminUsername: "azureuser",
				AdminPassword: "",
			},
		},
	}

	existingCS := &ContainerService{
		Properties: &Properties{
			ServicePrincipalProfile: &ServicePrincipalProfile{
				ClientID: "existingFakeID",
				Secret:   "existingSecret",
			},
			WindowsProfile: &WindowsProfile{
				AdminUsername: "azureuser",
				AdminPassword: "existingPassword",
			},
		},
	}
	if err := newCS.Merge(existingCS); err != nil {
		t.Fatalf("unexpectedly detected merge failure, %+v", err)
	}
	if newCS.Properties.ServicePrincipalProfile.ClientID != "fakeID" {
		t.Fatalf("unexpected Properties.ServicePrincipalProfile.ClientID changed")
	}
	if newCS.Properties.ServicePrincipalProfile.Secret != "existingSecret" {
		t.Fatalf("unexpected Properties.ServicePrincipalProfile.Secret not updated")
	}
	if newCS.Properties.WindowsProfile.AdminPassword != "existingPassword" {
		t.Fatalf("unexpected Properties.WindowsProfile.AdminPassword not updated")
	}
}
