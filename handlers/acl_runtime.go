package handlers

import (
	"github.com/go-openapi/runtime/middleware"
	client_native "github.com/haproxytech/client-native/v6"
	"github.com/haproxytech/client-native/v6/models"
	"github.com/haproxytech/dataplaneapi/misc"
	"github.com/haproxytech/dataplaneapi/operations/acl_runtime"
)

type GetACLSHandlerRuntimeImpl struct {
	Client client_native.HAProxyClient
}

func (h GetACLSHandlerRuntimeImpl) Handle(params acl_runtime.GetServicesHaproxyRuntimeAclsParams, i interface{}) middleware.Responder {
	runtime, err := h.Client.Runtime()
	if err != nil {
		e := misc.HandleError(err)
		return acl_runtime.NewGetServicesHaproxyRuntimeAclsDefault(int(*e.Code)).WithPayload(e)
	}

	files, err := runtime.GetACLFiles()
	if err != nil {
		e := misc.HandleError(err)
		return acl_runtime.NewGetServicesHaproxyRuntimeAclsDefault(int(*e.Code)).WithPayload(e)
	}

	return acl_runtime.NewGetServicesHaproxyRuntimeAclsOK().WithPayload(files)
}

type GetACLHandlerRuntimeImpl struct {
	Client client_native.HAProxyClient
}

func (h GetACLHandlerRuntimeImpl) Handle(params acl_runtime.GetServicesHaproxyRuntimeAclsIDParams, i interface{}) middleware.Responder {
	runtime, err := h.Client.Runtime()
	if err != nil {
		e := misc.HandleError(err)
		return acl_runtime.NewGetServicesHaproxyRuntimeAclsIDDefault(int(*e.Code)).WithPayload(e)
	}

	aclFile, err := runtime.GetACLFile(params.ID)
	if err != nil {
		e := misc.HandleError(err)
		return acl_runtime.NewGetServicesHaproxyRuntimeAclsIDDefault(int(*e.Code)).WithPayload(e)
	}

	return acl_runtime.NewGetServicesHaproxyRuntimeAclsIDOK().WithPayload(aclFile)
}

type GetACLFileEntriesHandlerRuntimeImpl struct {
	Client client_native.HAProxyClient
}

func (h GetACLFileEntriesHandlerRuntimeImpl) Handle(params acl_runtime.GetServicesHaproxyRuntimeAclsParentNameEntriesParams, i interface{}) middleware.Responder {
	runtime, err := h.Client.Runtime()
	if err != nil {
		e := misc.HandleError(err)
		return acl_runtime.NewGetServicesHaproxyRuntimeAclsParentNameEntriesDefault(int(*e.Code)).WithPayload(e)
	}
	files, err := runtime.GetACLFilesEntries(params.ParentName)
	if err != nil {
		e := misc.HandleError(err)
		return acl_runtime.NewGetServicesHaproxyRuntimeAclsParentNameEntriesDefault(int(*e.Code)).WithPayload(e)
	}

	return acl_runtime.NewGetServicesHaproxyRuntimeAclsParentNameEntriesOK().WithPayload(files)
}

type PostACLFileEntryHandlerRuntimeImpl struct {
	Client client_native.HAProxyClient
}

func (h PostACLFileEntryHandlerRuntimeImpl) Handle(params acl_runtime.PostServicesHaproxyRuntimeAclsParentNameEntriesParams, i interface{}) middleware.Responder {
	var err error
	runtime, err := h.Client.Runtime()
	if err != nil {
		e := misc.HandleError(err)
		return acl_runtime.NewPostServicesHaproxyRuntimeAclsParentNameEntriesDefault(int(*e.Code)).WithPayload(e)
	}

	if err = runtime.AddACLFileEntry(params.ParentName, params.Data.Value); err != nil {
		e := misc.HandleError(err)
		return acl_runtime.NewPostServicesHaproxyRuntimeAclsParentNameEntriesDefault(int(*e.Code)).WithPayload(e)
	}

	var fileEntry *models.ACLFileEntry
	fileEntry, err = runtime.GetACLFileEntry(params.ParentName, params.Data.Value)
	if err != nil {
		e := misc.HandleError(err)
		return acl_runtime.NewPostServicesHaproxyRuntimeAclsParentNameEntriesDefault(int(*e.Code)).WithPayload(e)
	}

	return acl_runtime.NewPostServicesHaproxyRuntimeAclsParentNameEntriesCreated().WithPayload(fileEntry)
}

type GetACLFileEntryRuntimeImpl struct {
	Client client_native.HAProxyClient
}

func (h GetACLFileEntryRuntimeImpl) Handle(params acl_runtime.GetServicesHaproxyRuntimeAclsParentNameEntriesIDParams, i interface{}) middleware.Responder {
	runtime, err := h.Client.Runtime()
	if err != nil {
		e := misc.HandleError(err)
		return acl_runtime.NewGetServicesHaproxyRuntimeAclsParentNameEntriesIDDefault(int(*e.Code)).WithPayload(e)
	}

	fileEntry, err := runtime.GetACLFileEntry(params.ParentName, params.ID)
	if err != nil {
		e := misc.HandleError(err)
		return acl_runtime.NewGetServicesHaproxyRuntimeAclsParentNameEntriesIDDefault(int(*e.Code)).WithPayload(e)
	}

	return acl_runtime.NewGetServicesHaproxyRuntimeAclsParentNameEntriesIDOK().WithPayload(fileEntry)
}

type DeleteACLFileEntryHandlerRuntimeImpl struct {
	Client client_native.HAProxyClient
}

func (h DeleteACLFileEntryHandlerRuntimeImpl) Handle(params acl_runtime.DeleteServicesHaproxyRuntimeAclsParentNameEntriesIDParams, i interface{}) middleware.Responder {
	runtime, err := h.Client.Runtime()
	if err != nil {
		e := misc.HandleError(err)
		return acl_runtime.NewDeleteServicesHaproxyRuntimeAclsParentNameEntriesIDDefault(int(*e.Code)).WithPayload(e)
	}
	if err := runtime.DeleteACLFileEntry(params.ParentName, "#"+params.ID); err != nil {
		e := misc.HandleError(err)
		return acl_runtime.NewDeleteServicesHaproxyRuntimeAclsParentNameEntriesIDDefault(int(*e.Code)).WithPayload(e)
	}

	return acl_runtime.NewDeleteServicesHaproxyRuntimeAclsParentNameEntriesIDNoContent()
}

type ACLRuntimeAddPayloadRuntimeACLHandlerImpl struct {
	Client client_native.HAProxyClient
}

func (h ACLRuntimeAddPayloadRuntimeACLHandlerImpl) Handle(params acl_runtime.AddPayloadRuntimeACLParams, i interface{}) middleware.Responder {
	runtime, err := h.Client.Runtime()
	if err != nil {
		e := misc.HandleError(err)
		return acl_runtime.NewAddPayloadRuntimeACLDefault(int(*e.Code)).WithPayload(e)
	}

	err = runtime.AddACLAtomic(params.ParentName, params.Data)
	if err != nil {
		status := misc.GetHTTPStatusFromErr(err)
		return acl_runtime.NewAddPayloadRuntimeACLDefault(status).WithPayload(misc.SetError(status, err.Error()))
	}
	return acl_runtime.NewAddPayloadRuntimeACLCreated().WithPayload(params.Data)
}
