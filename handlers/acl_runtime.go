package handlers

import (
	"github.com/go-openapi/runtime/middleware"
	client_native "github.com/haproxytech/client-native/v2"
	"github.com/haproxytech/client-native/v2/models"
	"github.com/haproxytech/dataplaneapi/misc"
	"github.com/haproxytech/dataplaneapi/operations/acl_runtime"
)

type GetACLSHandlerRuntimeImpl struct {
	Client *client_native.HAProxyClient
}

func (g GetACLSHandlerRuntimeImpl) Handle(params acl_runtime.GetServicesHaproxyRuntimeAclsParams, i interface{}) middleware.Responder {
	files, err := g.Client.Runtime.GetACLFiles()
	if err != nil {
		e := misc.HandleError(err)
		return acl_runtime.NewGetServicesHaproxyRuntimeAclsDefault(int(*e.Code)).WithPayload(e)
	}

	return acl_runtime.NewGetServicesHaproxyRuntimeAclsOK().WithPayload(files)
}

type GetACLHandlerRuntimeImpl struct {
	Client *client_native.HAProxyClient
}

func (g GetACLHandlerRuntimeImpl) Handle(params acl_runtime.GetServicesHaproxyRuntimeAclsIDParams, i interface{}) middleware.Responder {
	aclFile, err := g.Client.Runtime.GetACLFile(params.ID)
	if err != nil {
		e := misc.HandleError(err)
		return acl_runtime.NewGetServicesHaproxyRuntimeAclsIDDefault(int(*e.Code)).WithPayload(e)
	}

	return acl_runtime.NewGetServicesHaproxyRuntimeAclsIDOK().WithPayload(aclFile)
}

type GetACLFileEntriesHandlerRuntimeImpl struct {
	Client *client_native.HAProxyClient
}

func (g GetACLFileEntriesHandlerRuntimeImpl) Handle(params acl_runtime.GetServicesHaproxyRuntimeACLFileEntriesParams, i interface{}) middleware.Responder {
	files, err := g.Client.Runtime.GetACLFilesEntries(params.ACLID)
	if err != nil {
		e := misc.HandleError(err)
		return acl_runtime.NewGetServicesHaproxyRuntimeACLFileEntriesDefault(int(*e.Code)).WithPayload(e)
	}

	return acl_runtime.NewGetServicesHaproxyRuntimeACLFileEntriesOK().WithPayload(files)
}

type PostACLFileEntryHandlerRuntimeImpl struct {
	Client *client_native.HAProxyClient
}

func (c PostACLFileEntryHandlerRuntimeImpl) Handle(params acl_runtime.PostServicesHaproxyRuntimeACLFileEntriesParams, i interface{}) middleware.Responder {
	var err error

	if err = c.Client.Runtime.AddACLFileEntry(params.ACLID, params.Data.Value); err != nil {
		e := misc.HandleError(err)
		return acl_runtime.NewPostServicesHaproxyRuntimeACLFileEntriesDefault(int(*e.Code)).WithPayload(e)
	}

	var fileEntry *models.ACLFileEntry
	fileEntry, err = c.Client.Runtime.GetACLFileEntry(params.ACLID, params.Data.Value)
	if err != nil {
		e := misc.HandleError(err)
		return acl_runtime.NewPostServicesHaproxyRuntimeACLFileEntriesDefault(int(*e.Code)).WithPayload(e)
	}

	return acl_runtime.NewPostServicesHaproxyRuntimeACLFileEntriesCreated().WithPayload(fileEntry)
}

type GetACLFileEntryRuntimeImpl struct {
	Client *client_native.HAProxyClient
}

func (g GetACLFileEntryRuntimeImpl) Handle(params acl_runtime.GetServicesHaproxyRuntimeACLFileEntriesIDParams, i interface{}) middleware.Responder {
	fileEntry, err := g.Client.Runtime.GetACLFileEntry(params.ACLID, params.ID)
	if err != nil {
		e := misc.HandleError(err)
		return acl_runtime.NewGetServicesHaproxyRuntimeACLFileEntriesIDDefault(int(*e.Code)).WithPayload(e)
	}

	return acl_runtime.NewGetServicesHaproxyRuntimeACLFileEntriesIDOK().WithPayload(fileEntry)
}

type DeleteACLFileEntryHandlerRuntimeImpl struct {
	Client *client_native.HAProxyClient
}

func (d DeleteACLFileEntryHandlerRuntimeImpl) Handle(params acl_runtime.DeleteServicesHaproxyRuntimeACLFileEntriesIDParams, i interface{}) middleware.Responder {
	if err := d.Client.Runtime.DeleteACLFileEntry(params.ACLID, "#"+params.ID); err != nil {
		e := misc.HandleError(err)
		return acl_runtime.NewDeleteServicesHaproxyRuntimeACLFileEntriesIDDefault(int(*e.Code)).WithPayload(e)
	}

	return acl_runtime.NewDeleteServicesHaproxyRuntimeACLFileEntriesIDNoContent()
}
