// var rolesett = kendo.observable({
//     titleModel: "New Roles",
//     loading: false,
//     edit: true,
//     TitleFilter: " Hide Filter",
//     disableRolename: true,
//     temp: [],
//     filterStatus: "",
//     Id: "",
//     roleName: "",
//     page: ""
// })
var rolesett = {
    titleModel: ko.observable("New Roles"),
    loading: ko.observable(false),
    edit: ko.observable(true),
    TitleFilter: ko.observable(" Hide Filter"),
    disableRolename: ko.observable(true),
    temp: ko.observableArray([]),
    filterStatus: ko.observable(""),
    Id: ko.observable(""),
    roleName: ko.observable(""),
    page: ko.observable(""),
    DatePageBar : ko.observable(),
}

rolesett.mappingRole={
    name: ko.observable(""),
    status: ko.observable(false),
    menu: ko.observableArray([]),
    landing:ko.observable(""),
}
// rolesett.mappingRole = kendo.observable({
//     name: "",
//     status: false,
//     menu: [],
//     landing: ""
// })
var msFilterRoleVM = {
    filterRole: ko.observableArray([]),
    listRole: ko.observableArray([]),
}
// var msFilterRoleVM = kendo.observable({
//     filterRole: [],
//     listRole: []
// })
var ddRoleVM ={
    filterPage: ko.observableArray([]),
    listPage: ko.observableArray([]),
}
// var ddRoleVM = kendo.observable({
//     filterPage: [],
//     listPage: []
// })

var masterGridMenuData = new kendo.data.ObservableArray([]);

// rolesett.notLoading = function(){
//     // return ( ! rolesett.get("loading"))
//     return ( ! rolesett.loading())
// }
var notLoading = ko.computed(
    function(){
        return ( ! rolesett.loading())
    }
)

rolesett.getDateNow = function(){
    var page = moment(new Date()).format("DD MMMM YYYY HH:mm")
    rolesett.DatePageBar(page)
}


rolesett.clickSave = function(){
    if(rolesett.edit() != true)
    // if(rolesett.get("edit") != true)
        rolesett.SaveData()
    else
        rolesett.UpdateData()
}

rolesett.ClearField = function(){
    // rolesett.set("roleName", "");
    rolesett.roleName("");
    $('#Status').bootstrapSwitch('state',true);
}

rolesett.Search = function(){
    rolesett.GetDataRole();
}

rolesett.Reset = function(){
    var multiSelect = $('#filterRole').data("kendoMultiSelect");
    multiSelect.value([]);
    $('#filterStatus').bootstrapSwitch('state', true);
    rolesett.GetDataRole();
}

rolesett.AddNew = function(){
    var landing = $("#role").data("kendoDropDownList");
    landing.value("");
    $("#roleModal").modal("show");
    $("#nav-dex").css('z-index', '0');
    $("#roleModal").modal({
        backdrop: 'static',
        keyboard: false
    });
    $('.rolecheck-value-check-all').prop('checked', false);
    $('.rolecheck-value-Access').prop('checked', false);
    $('.rolecheck-value-Create').prop('checked', false);
    $('.rolecheck-value-Delete').prop('checked', false);
    $('.rolecheck-value-Edit').prop('checked', false);
    $('.rolecheck-value-View').prop('checked', false);
    $('.rolecheck-value-Approve').prop('checked', false);
    $('.rolecheck-value-Process').prop('checked', false);
    // rolesett.set("titleModel", "New Roles");
    rolesett.titleModel("New Roles")
    rolesett.disableRolename(true)
    // rolesett.set("disableRolename", true);
    rolesett.ClearField();
    rolesett.edit(true)
    // rolesett.set("edit", true);
    rolesett.getTopMenu();
}

rolesett.SaveData = function(){
    var displayedData = $("#MasterGridMenu").data().kendoTreeList.dataSource.view();
    rolesett.mappingRole.name(rolesett.roleName())
    rolesett.mappingRole.status($('#Status').bootstrapSwitch('state'))
    rolesett.mappingRole.landing($('#role').val()) 
    // rolesett.mappingRole.set("name", rolesett.get("roleName"));
    // rolesett.mappingRole.set("status", $('#Status').bootstrapSwitch('state'));
    // rolesett.mappingRole.set("landing", $('#role').val());
    
    // rolesett.mappingRole.set("menu", []);
    rolesett.mappingRole.menu([]);
    for (var i in displayedData){
        if (displayedData[i].Id != undefined){
            var Access = $("#check-Access-"+displayedData[i].Id).is(":checked");
            // if (Access != false){
                rolesett.mappingRole.menu.push({
                    "menuid" : displayedData[i].Id,
                    "menuname" : displayedData[i].Title,
                    "haschild" : displayedData[i].Haschild,
                    "enable" : displayedData[i].Enable,
                    "parentid" : displayedData[i].Parent,
                    "checkall" : $("#check-all-new"+displayedData[i].Id).is(":checked"),
                    "access" : $("#check-Access-"+displayedData[i].Id).is(":checked"),
                    "view" : $("#check-View-"+displayedData[i].Id).is(":checked"),
                    "create" : $("#check-Create-"+displayedData[i].Id).is(":checked"),
                    "approve" : $("#check-Approve-"+displayedData[i].Id).is(":checked"),
                    "delete" : $("#check-Delete-"+displayedData[i].Id).is(":checked"),
                    "process" : $("#check-Process-"+displayedData[i].Id).is(":checked"),
                    "edit" : $("#check-Edit-"+displayedData[i].Id).is(":checked"),
                    "Url": displayedData[i].Url
                });
            // }
        }


    }
    

    var param = rolesett.mappingRole;
    console.log("param",param);
    var url = "/sysroles/savedata";
    var validator = $("#AddRole").data("kendoValidator");
    if(validator == undefined){
       validator= $("#AddRole").kendoValidator().data("kendoValidator");
    }
    // rolesett.Cancel();
    // rolesett.Reset();
    if (validator.validate()) {
        ajaxPost(url, param, function(res){
            if(res.IsError != true){
                rolesett.Cancel();
                rolesett.Reset();
                $("#nav-dex").css('z-index', 'none');
                swal({
                    title:"Success!",
                    text: res.Message,
                    type: "success",
                    confirmButtonColor:"#3da09a"});
                location.reload();
            }else{
                return swal({
                    title:"Error!",
                    text: res.Message,
                    type: "error",
                    confirmButtonColor:"#3da09a"
                });
            }
        });
    }
}

rolesett.UpdateData = function(){
    var displayedData = $("#MasterGridMenu").data().kendoTreeList.dataSource.view();
    console.log(displayedData)
    rolesett.mappingRole.name(rolesett.roleName())
    rolesett.mappingRole.status($('#Status').bootstrapSwitch('state'))
    rolesett.mappingRole.landing($('#role').val()) 
    // rolesett.mappingRole.set("name", rolesett.get("roleName"));
    // rolesett.mappingRole.set("status", $('#Status').bootstrapSwitch('state'));
    // rolesett.mappingRole.set("landing", $('#role').val());
    
    // rolesett.mappingRole.set("menu", []);
    rolesett.mappingRole.menu([]);
    for (var i in displayedData){
        if (displayedData[i].Id != undefined){
            var Access = $("#check-Access-"+displayedData[i].Id).is(":checked");
            // if (Access != false){
                rolesett.mappingRole.menu.push({
                    "menuid" : displayedData[i].Id,
                    "menuname" : displayedData[i].Title,
                    "haschild" : displayedData[i].Haschild,
                    "enable" : displayedData[i].Enable,
                    "parentid" : displayedData[i].Parent,
                    "checkall" : $("#check-all-new"+displayedData[i].Id).is(":checked"),
                    "access" : $("#check-Access-"+displayedData[i].Id).is(":checked"),
                    "view" : $("#check-View-"+displayedData[i].Id).is(":checked"),
                    "create" : $("#check-Create-"+displayedData[i].Id).is(":checked"),
                    "approve" : $("#check-Approve-"+displayedData[i].Id).is(":checked"),
                    "delete" : $("#check-Delete-"+displayedData[i].Id).is(":checked"),
                    "process" : $("#check-Process-"+displayedData[i].Id).is(":checked"),
                    "edit" : $("#check-Edit-"+displayedData[i].Id).is(":checked"),
                    "Url": displayedData[i].Url
                });
            // }
        }
    }

    var param =  rolesett.mappingRole;
    console.log('param',param)
    param.Id = rolesett.Id();
    // param.Id = rolesett.get("Id");
    var url = "/sysroles/savedata";
    var validator = $("#AddRole").data("kendoValidator");
    if(validator==undefined){
       validator= $("#AddRole").kendoValidator().data("kendoValidator");
    }
    if (validator.validate()) {
        ajaxPost(url, param, function(res){
            if(res.IsError != true){
                rolesett.Cancel();
                rolesett.Reset();
                $("#nav-dex").css('z-index', 'none');
                swal({
                    title:"Success!",
                    text: res.Message,
                    type: "success",
                    confirmButtonColor:"#3da09a"});
                location.reload();
            }else{
                return swal({
                    title:"Error!",
                    text: res.Message,
                    type: "error",
                    confirmButtonColor:"#3da09a"});
            }
        });
    }
}

rolesett.EditData = function(IdRole){
    var url = "/sysroles/getmenuedit";
    var param = {
            Id : IdRole
        }
        ajaxPost(url, param, function(res){
            if(res.IsError != true){
                rolesett.disableRolename(false);
                // rolesett.set("disableRolename", false);
                $("#roleModal").modal("show");
                $("#nav-dex").css('z-index', '0');
                $("#roleModal").modal({
                    backdrop: 'static',
                    keyboard: false
                });
                rolesett.titleModel("Update Roles");
                rolesett.edit(true);
                // rolesett.set("titleModel", "Update Roles");
                // rolesett.set("edit", true);
                var Records = res.Data.Records[0];
                rolesett.Id(Records.Id);
                // rolesett.set("Id", Records.Id);
                rolesett.roleName(Records.Name);
                // rolesett.set("roleName", Records.roleName);
                var landing = $("#role").data("kendoDropDownList");
                landing.value(Records.Landing);
                $('#Status').bootstrapSwitch('state',Records.Status);
                var dataMenu = res.Data.Records[0].Menu;
                var newRecords = [];
                for (var d in dataMenu){
                    newRecords.push({
                        "Access": dataMenu[d].Access,
                        "Approve": dataMenu[d].Approve,
                        "Checkall": dataMenu[d].Checkall,
                        "Create": dataMenu[d].Create,
                        "Delete": dataMenu[d].Delete,
                        "Edit": dataMenu[d].Edit,
                        "Enable": dataMenu[d].Enable,
                        "Haschild": dataMenu[d].Haschild,
                        "Id": dataMenu[d].Menuid,
                        "Title": dataMenu[d].Menuname,
                        "Parent": dataMenu[d].Parent,
                        "Process": dataMenu[d].Process,
                        "Url": dataMenu[d].Url,
                        "View": dataMenu[d].View,
                    });
                }
                rolesett.GetDataMenu(newRecords);
            }else{
                return swal({
                    title:"Error!",
                    text: res.Message,
                    type: "error",
                    confirmButtonColor:"#3da09a"});
            }
    }); 
}

rolesett.Cancel = function(){
    $("#roleModal").modal("hide");
    $("#nav-dex").css('z-index', 'none');
}

// rolesett.bind("change", function(e) {
//     if(e.field = "filterRole")
//       if(userinfo.get("view") != "false")
//         rolesett.GetDataRole();
// });

// var userid = userinfo.get("usernameh");
var userid = userinfo.usernameh()
var gc = new GridColumn('role_master', userid, 'MasterGridRole');

rolesett.GetDataRole = function(){
    // rolesett.set("loading", false);
    rolesett.loading(false);
    var param =  {
        // "Name" : msFilterRoleVM.get("filterRole"),
        "Name" : msFilterRoleVM.filterRole(),
        "Status" : $('#filterStatus').bootstrapSwitch('state')
    };
    var dataSource = [];
    var url = "/sysroles/getdata";
    $("#MasterGridRole").html("");
    $("#MasterGridRole").kendoGrid({
            dataSource: {
                    transport: {
                        read: {
                            url: url,
                            data: param,
                            dataType: "json",
                            type: "POST",
                            contentType: "application/json",
                        },
                        parameterMap: function(data) {                                 
                           return JSON.stringify(data);                                 
                        },
                    },
                    schema: {
                        data: function(data) {
                            gc.Init();
                            // rolesett.set("loading", false);
                            rolesett.loading(false);
                            if (data.Data.Count == 0) {
                                return dataSource;
                            } else {
                                return data.Data.Records;
                            }
                        },
                        total: function(data){
                            if (data.Data.Count == 0) {
                                return 0;
                            } else {
                                return data.Data.Records.length;
                            }
                        },
                    },
                    // pageSize: 15,
                    // serverPaging: true,
                    // serverSorting: true,
                },
                resizable: true,
                sortable: true,
                // pageable: {
                //     refresh: true,
                //     pageSizes: true,
                //     buttonCount: 5
                // },
                columnMenu: false,
                  columnHide: function(e) {
                    gc.RemoveColumn(e.column.field);
                  },
                  columnShow: function(e) {
                    gc.AddColumn(e.column.field);
                  },
            columns: [
                {
                    field:"Name",
                    title:"Role Name",
                    // width:150,
                    headerAttributes: {class: 'k-header header-bgcolor'},
                    template: "#if(rolesett.edit() != 'false'){#<a class='onclickusersetting' id='ls' href='javascript:rolesett.EditData(\"#: Id #\")'>#: Name #</a>#}else{#<div>#: Name #</div>#}#"
                    // template: "#if(rolesett.get('Edit') != 'false'){#<a class='grid-select' id='ls' href='javascript:rolesett.EditData(\"#: Id #\")'>#: Name #</a>#}else{#<div>#: Name #</div>#}#"
                },
                {
                    field:"Status",
                    title:"Status",
                    headerAttributes: {class: 'k-header header-bgcolor'},
                    template: function(res){
                        if(res.Status == true){
                            return "Active"
                        }else{
                            return "Not Active"
                        }
                    }
                    // width:50

                },
                ]
    });
}

rolesett.getTopMenu = function(){
    var param = {
    }
    var url = "/sysroles/getmenu";
    ajaxPost(url, param, function(res){
        if(res.IsError != true){
            var dataMenu = res.Data.Records;
            var newRecords = [];
            for (var d in dataMenu){
                newRecords.push({
                    "Access": false,
                    "Approve":false,
                    "Checkall": false,
                    "Create": false,
                    "Delete": false,
                    "Edit": false,
                    "Enable": true,
                    "Haschild": false,
                    "Id": dataMenu[d].Id,
                    "Title": dataMenu[d].Title,
                    "Parent": dataMenu[d].Parent,
                    "Process": false,
                    "Url": dataMenu[d].Url,
                    "View": false,
                    
                });
            }
            rolesett.GetDataMenu(newRecords);
        }else{
            return swal({
                    title:"Error!",
                    text: res.Message,
                    type: "error",
                    confirmButtonColor:"#3da09a"});
        }
    });
}

rolesett.GetDataMenu = function(e){
    $("#MasterGridMenu").data("kendoTreeList").setDataSource(e);
}

rolesett.unCheck = function(Menuid){
    if(!$("#check-Access-"+Menuid).prop('checked') || !$("#check-Create-"+Menuid).prop('checked') || !$("#check-Edit-"+Menuid).prop('checked') || !$("#check-Delete-"+Menuid).prop('checked') || !$("#check-View-"+Menuid).prop('checked') || !$("#check-Approve-"+Menuid).prop('checked') || !$("#check-Process-"+Menuid).prop('checked')){
        $('#check-all'+Menuid).prop('checked', false);
        $('#check-all-new'+Menuid).prop('checked', false);
    }else if($("#check-Access-"+Menuid).prop('checked') == true && $("#check-Create-"+Menuid).prop('checked')== true && $("#check-Edit-"+Menuid).prop('checked') == true && $("#check-Delete-"+Menuid).prop('checked') == true && $("#check-View-"+Menuid).prop('checked') == true && $("#check-Approve-"+Menuid).prop('checked') == true && $("#check-Process-"+Menuid).prop('checked') == true){
        $('#check-all'+Menuid).prop('checked', true);
        $('#check-all-new'+Menuid).prop('checked', true); 
    }
}

rolesett.Checkall = function(Menuid){
    $('#check-all'+Menuid).change(function(){
        $("#check-Access-"+Menuid).prop('checked', $(this).prop('checked'));
        $("#check-Create-"+Menuid).prop('checked', $(this).prop('checked'));
        $("#check-Edit-"+Menuid).prop('checked', $(this).prop('checked'));
        $("#check-Delete-"+Menuid).prop('checked', $(this).prop('checked'));
        $("#check-View-"+Menuid).prop('checked', $(this).prop('checked'));
        $("#check-Approve-"+Menuid).prop('checked', $(this).prop('checked'));
        $("#check-Process-"+Menuid).prop('checked', $(this).prop('checked'));
    });
    $('#check-all-new'+Menuid).change(function(){
        $("#check-Access-"+Menuid).prop('checked', $(this).prop('checked'));
        $("#check-Create-"+Menuid).prop('checked', $(this).prop('checked'));
        $("#check-Edit-"+Menuid).prop('checked', $(this).prop('checked'));
        $("#check-Delete-"+Menuid).prop('checked', $(this).prop('checked'));
        $("#check-View-"+Menuid).prop('checked', $(this).prop('checked'));
        $("#check-Approve-"+Menuid).prop('checked', $(this).prop('checked'));
        $("#check-Process-"+Menuid).prop('checked', $(this).prop('checked'));
    });
}

rolesett.getRole = function(){
    var param = {
    }
    var url = "/datamaster/getroles";
    msFilterRoleVM.listRole([]);
    // msFilterRoleVM.set("listRole", []);
    ajaxPost(url, param, function(res){
        var dataRole = Enumerable.From(res).OrderBy("$.name").ToArray();
        for (var r in dataRole){
            msFilterRoleVM.listRole.push({
                "text" : dataRole[r].name,
                "value" : dataRole[r].name,
            });
        }
    });
}

rolesett.toggleFilter = function(){
  var panelFilter = $('.panel-filter');
  var panelContent = $('.panel-content');

  if (panelFilter.is(':visible')) {
    panelFilter.hide();
    panelContent.attr('class', 'col-md-12 col-sm-12 ez panel-content');
    $('.breakdown-filter').removeAttr('style');
  } else {
    panelFilter.show();
    panelContent.attr('class', 'col-md-9 col-sm-9 ez panel-content');
    //panelContent.css('margin-top', '1.3%');
    $('.breakdown-filter').css('width', '60%');
  }

  $('.k-grid').each(function (i, d) {
    try {
      $(d).data('kendoGrid').refresh();
    } catch (err) {}
  });

  $('.k-pivot').each(function (i, d) {
    $(d).data('kendoPivotGrid').refresh();
  });

  $('.k-chart').each(function (i, d) {
    $(d).data('kendoChart').redraw();
  });
  rolesett.panel_relocated();
  // var FilterTitle = rolesett.get("TitleFilter");
  var FilterTitle = rolesett.TitleFilter();
  if (FilterTitle == " Hide Filter") {
        rolesett.TitleFilter(" Show Filter");
    } else if(FilterTitle == " Show Filter"){
        rolesett.TitleFilter(" Hide Filter");
    }
}

rolesett.panel_relocated = function(){
  if ($('.panel-yo').size() == 0) {
    return;
  }

  var window_top = $(window).scrollTop();
  var div_top = $('.panel-yo').offset().top;
  if (window_top > div_top) {
    $('.panel-fix').css('width', $('.panel-yo').width());
    $('.panel-fix').addClass('contentfilter');
    $('.panel-yo').height($('.panel-fix').outerHeight());
  } else {
    $('.panel-fix').removeClass('contentfilter');
    $('.panel-yo').height(0);
  }
}

rolesett.getLandingPage = function(){
    var param = {
    }
    var url = "/sysroles/getlandingpage";
    // ddRoleVM.set("listPage", []);
    ddRoleVM.listPage([]);
    ajaxPost(url, param, function(res){
        var dataPage = Enumerable.From(res).OrderBy("$.Title").ToArray();
        for (var u in dataPage){
            ddRoleVM.listPage.push({
                "text" : dataPage[u].Title,
                "value" : dataPage[u].Title,
            });
        }
        rolesett.DropDownLandingPage()
    });
}
rolesett.DropDownLandingPage = function(){
    console.log(ddRoleVM.listPage())
    $('#role').kendoDropDownList({
        dataSource : ddRoleVM.listPage(),
        valuePrimitive: true,
        dataTextField: 'text', 
        dataValueField: 'value', 
        optionLabel: 'Select landing page...',
        change: function(e){
            var key = this.dataItem();
            // ddRoleVM.filterPage(key.value)
        }
    })
}
$(document).ready(function (){ 
    // kendo.bind($("#rolesett"), rolesett)
    // kendo.bind($("#filterRole"), msFilterRoleVM)
    // kendo.bind($("#role"), ddRoleVM)
    $("#MasterGridMenu").kendoTreeList({
        dataSource: {
            data: [],
            schema: {
                model: {
                    id: "Id",
                    parentId: "Parent",
                    fields: {
                        _id: { field: "_id", type: "string" },
                        titleFolder: { field: "titleFolder", type: "string" },
                        parentId: { field: "Parent", type: "string" }
                    },
                    expanded: true
                }
            }
        },
        height: 400,
        columns: [
            { 
                field: "Title",
                title:"Title", 
                width: 200 
            },
            { 
                field:"Checkall",
                title:"Check All", 
                width: 50,
                attributes:{"class": "align-center"},
                template: "#if(parentId != '' || Haschild == false){#<input id='check-all-new#:Id#' class='rolecheck-value-check-all' type='checkbox' onclick='rolesett.Checkall(#:Id#)' #: Checkall==true ? 'checked' : '' #/>#}#"
            },
            {
                field:"Access",
                title:"Access",
                width:50,
                attributes: {"class": "align-center"},
                template:"#if(parentId != '' || Haschild == false){#<input id='check-Access-#:Id #' class='rolecheck-value-Access' onclick='rolesett.unCheck(#:Id#)' type='checkbox' #: Access==true ? 'checked' : '' #/>#}#"              
            },
            {
                field:"Create",
                title:"Create",
                width:50,
                attributes: {"class": "align-center"},
                template:"#if(parentId != '' || Haschild == false){#<input id='check-Create-#:Id #' class='rolecheck-value-Create' onclick='rolesett.unCheck(#:Id#)' type='checkbox' #: Create==true ? 'checked' : '' #/>#}#"              
            },
            {
                field:"Edit",
                title:"Edit",
                width:50,
                attributes: {"class": "align-center"},
                template:"#if(parentId != '' || Haschild == false){#<input id='check-Edit-#:Id #' class='rolecheck-value-Edit' onclick='rolesett.unCheck(#:Id#)' type='checkbox' #: Edit==true ? 'checked' : '' #/>#}#"  
            },
            {
                field:"Delete",
                title:"Delete",
                width:50,
                attributes: {"class": "align-center"},
                template:"#if(parentId != '' || Haschild == false){#<input id='check-Delete-#:Id #' class='rolecheck-value-Delete' onclick='rolesett.unCheck(#:Id#)' type='checkbox' #: Delete==true ? 'checked' : '' #/>#}#"
            },
            {
                field:"View",
                title:"View",
                width:50,
                attributes: {"class": "align-center"},
                template:"#if(parentId != '' || Haschild == false){#<input id='check-View-#:Id #' class='rolecheck-value-View' onclick='rolesett.unCheck(#:Id#)' type='checkbox' #: View==true ? 'checked' : '' #/>#}#"
            },
            {
                field:"Approve",
                title:"Approve",
                width:50,
                attributes: {"class": "align-center"},
                template:"#if(parentId != '' || Haschild == false){#<input id='check-Approve-#:Id #' class='rolecheck-value-Approve' onclick='rolesett.unCheck(#:Id#)' type='checkbox' #: Approve==true ? 'checked' : '' #/>#}#"
            },
            {
                field:"Process",
                title:"Process",
                width:50,
                attributes: {"class": "align-center"},
                template:"#if(parentId != '' || Haschild == false){#<input id='check-Process-#:Id #' class='rolecheck-value-Process' onclick='rolesett.unCheck(#:Id#)' type='checkbox' #: Process==true ? 'checked' : '' #/>#}#"
            }
        ],
    }); 

    $('#filterStatus').bootstrapSwitch('state',true)
    rolesett.getRole();
    rolesett.GetDataRole();
    rolesett.getLandingPage();
    rolesett.getDateNow()
    // $("#filterRole").kendoMultiSelect({
    //     // dataSource: msFilterRoleVM.listRole(),
    //     dataValueField: 'value', 
    //     dataTextField: 'text',
    //     valuePrimitive: true,
    //     optionLabel: "Select role..."
    // })
    $(".k-icon.k-clear-value.k-i-close").hide()
});

// var roleGridVM = kendo.observable({
//     dataSource: new kendo.data.DataSource({
//     	transport: {
//             read: {
//                 url: "#",
//                 dataType: "json"
//             }
//         },
//         schema: {
//             model: {
//                 id: "ProductID",
//                 fields: {
//                     rolename: { type: "string" },
//                     status: { type: "string" },
//                 }
//             }
//         },
//     })
// });

// var grid = $('#MasterGridRole').kendoGrid({
//     dataSource: roleGridVM.dataSource,
//     columns: [
//      	{ 'field': 'rolename' }, 
//      	{ 'field': 'status' },
//     ],
// }).data('kendoGrid');

// kendo.bind($("#MasterGridUser"), roleGridVM);