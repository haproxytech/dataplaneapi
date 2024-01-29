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

func (h GetACLFileEntriesHandlerRuntimeImpl) Handle(params acl_runtime.GetServicesHaproxyRuntimeACLFileEntriesParams, i interface{}) middleware.Responder {
	runtime, err := h.Client.Runtime()
	if err != nil {
		e := misc.HandleError(err)
		return acl_runtime.NewGetServicesHaproxyRuntimeACLFileEntriesDefault(int(*e.Code)).WithPayload(e)
	}
	files, err := runtime.GetACLFilesEntries(params.ACLID)
	if err != nil {
		e := misc.HandleError(err)
		return acl_runtime.NewGetServicesHaproxyRuntimeACLFileEntriesDefault(int(*e.Code)).WithPayload(e)
	}

	return acl_runtime.NewGetServicesHaproxyRuntimeACLFileEntriesOK().WithPayload(files)
}

type PostACLFileEntryHandlerRuntimeImpl struct {
	Client client_native.HAProxyClient
}

func (h PostACLFileEntryHandlerRuntimeImpl) Handle(params acl_runtime.PostServicesHaproxyRuntimeACLFileEntriesParams, i interface{}) middleware.Responder {
	var err error
	runtime, err := h.Client.Runtime()
	if err != nil {
		e := misc.HandleError(err)
		return acl_runtime.NewPostServicesHaproxyRuntimeACLFileEntriesDefault(int(*e.Code)).WithPayload(e)
	}

	if err = runtime.AddACLFileEntry(params.ACLID, params.Data.Value); err != nil {
		e := misc.HandleError(err)
		return acl_runtime.NewPostServicesHaproxyRuntimeACLFileEntriesDefault(int(*e.Code)).WithPayload(e)
	}

	var fileEntry *models.ACLFileEntry
	fileEntry, err = runtime.GetACLFileEntry(params.ACLID, params.Data.Value)
	if err != nil {
		e := misc.HandleError(err)
		return acl_runtime.NewPostServicesHaproxyRuntimeACLFileEntriesDefault(int(*e.Code)).WithPayload(e)
	}

	return acl_runtime.NewPostServicesHaproxyRuntimeACLFileEntriesCreated().WithPayload(fileEntry)
}

type GetACLFileEntryRuntimeImpl struct {
	Client client_native.HAProxyClient
}

func (h GetACLFileEntryRuntimeImpl) Handle(params acl_runtime.GetServicesHaproxyRuntimeACLFileEntriesIDParams, i interface{}) middleware.Responder {
	runtime, err := h.Client.Runtime()
	if err != nil {
		e := misc.HandleError(err)
		return acl_runtime.NewGetServicesHaproxyRuntimeACLFileEntriesIDDefault(int(*e.Code)).WithPayload(e)
	}

	fileEntry, err := runtime.GetACLFileEntry(params.ACLID, params.ID)
	if err != nil {
		e := misc.HandleError(err)
		return acl_runtime.NewGetServicesHaproxyRuntimeACLFileEntriesIDDefault(int(*e.Code)).WithPayload(e)
	}

	return acl_runtime.NewGetServicesHaproxyRuntimeACLFileEntriesIDOK().WithPayload(fileEntry)
}

type DeleteACLFileEntryHandlerRuntimeImpl struct {
	Client client_native.HAProxyClient
}

func (h DeleteACLFileEntryHandlerRuntimeImpl) Handle(params acl_runtime.DeleteServicesHaproxyRuntimeACLFileEntriesIDParams, i interface{}) middleware.Responder {
	runtime, err := h.Client.Runtime()
	if err != nil {
		e := misc.HandleError(err)
		return acl_runtime.NewDeleteServicesHaproxyRuntimeACLFileEntriesIDDefault(int(*e.Code)).WithPayload(e)
	}
	if err := runtime.DeleteACLFileEntry(params.ACLID, "#"+params.ID); err != nil {
		e := misc.HandleError(err)
		return acl_runtime.NewDeleteServicesHaproxyRuntimeACLFileEntriesIDDefault(int(*e.Code)).WithPayload(e)
	}

	return acl_runtime.NewDeleteServicesHaproxyRuntimeACLFileEntriesIDNoContent()
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

	err = runtime.AddACLAtomic(params.ACLID, params.Data)
	if err != nil {
		status := misc.GetHTTPStatusFromErr(err)
		return acl_runtime.NewAddPayloadRuntimeACLDefault(status).WithPayload(misc.SetError(status, err.Error()))
	}
	return acl_runtime.NewAddPayloadRuntimeACLCreated().WithPayload(params.Data)
}
