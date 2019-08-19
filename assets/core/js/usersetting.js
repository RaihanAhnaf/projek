// var usersett = kendo.observable({
//     titleModel: "New User",
//     loading: false,
//     edit: false,
//     TitelFilter: " Hide Filter",
//     Id: "",
//     userName: "",
//     fullName: "",
//     email: "",
//     password: "",
//     confirmPassword: "",
// });
var usersett = {
    titleModel: ko.observable("New User"),
    loading: ko.observable(false),
    edit: ko.observable(false),
    TitelFilter: ko.observable(" Hide Filter"),
    Id: ko.observable(""),
    userName: ko.observable(""),
    fullName: ko.observable(""),
    email: ko.observable(""),
    password: ko.observable(""),
    potition: ko.observable(""),
    confirmPassword: ko.observable(""),
    DatePageBar: ko.observable(),
};

// var msUsernameVM = kendo.observable({
//     filterUser: [],
//     listUserName: [],
// });
var msUsernameVM = {
    filterUser: ko.observableArray([]),
    listUserName: ko.observableArray([]),
};

// var msRoleVM = kendo.observable({
//     filterRole: [],
//     role: "",
//     listRole: [],
// });
var msRoleVM = {
    filterRole: ko.observableArray([]),
    role: ko.observable(""),
    listRole: ko.observableArray([]),
};
var msLocation = {
    data: ko.observableArray([]),
    value: ko.observable(""),
}
// usersett.notLoading = function(){
//     return ( ! usersett.get("loading"))
// }

// usersett.notEdit = function(){
//     return ( ! usersett.get("edit"))
// }

usersett.clickSave = function() {
    // if(usersett.get("edit") != true)
    if (usersett.edit() != true)
        usersett.SaveData()
    else
        usersett.UpdateData()
}

usersett.getDateNow = function(){
    var page = moment(new Date()).format("DD MMMM YYYY HH:mm")
    usersett.DatePageBar(page)
}

usersett.ClearField = function() {
    usersett.Id("");
    usersett.userName("");
    usersett.fullName("");
    usersett.email("");
    usersett.password("");
    usersett.confirmPassword("");
    msRoleVM.role("")
    usersett.potition("")
    msLocation.value("")
    // usersett.set("Id", "");
    // usersett.set("userName", "");
    // usersett.set("fullName", "");
    // usersett.set("email", "");
    // usersett.set("password", "");
    // usersett.set("confirmPassword", "");
    // msRoleVM.set("role", "");
    $('#Status').bootstrapSwitch('state', true);
}

// usersett.Search = function() {
//     usersett.GetDataUser();
// }

usersett.Reset = function() {
    // msUsernameVM.set("filterUser", []);
    // msRoleVM.set("filterRole", []);
    msUsernameVM.filterUser([]);
    msRoleVM.filterRole([]);
    usersett.GetDataUser();
    $('#StatusFilter').bootstrapSwitch('state', true);
}

usersett.AddNew = function() {
    usersett.ClearField();
    $("#userModal").modal("show");
    $("#nav-dex").css('z-index', '0');
    $("#userModal").modal({
        backdrop: 'static',
        keyboard: false
    });

    // usersett.set("titleModel", "New User");
    usersett.titleModel("New User");

    $('#filterStatus').bootstrapSwitch('state', true)
    usersett.edit(false)
    // usersett.set("edit", false);

    var validator = $("#AddUserSetting").kendoValidator().data("kendoValidator");
    validator.hideMessages();
}

usersett.SaveData = function() {
    var locationId = parseInt(msLocation.value())
    var location = _.find(msLocation.data(), function(e){return e.LocationID == locationId})
    var statusBool = $('#Status').bootstrapSwitch('state');
    var dropRole = $("#role").data("kendoDropDownList");
    // var param = {
    //     "UserName": usersett.get("userName"),
    //     "FullName": usersett.get("fullName"),
    //     "Email": usersett.get("email"),
    //     "Enable": statusBool,
    //     "Password": usersett.get("password"),
    //     "Role": dropRole.text(),
    // }
    var param = {
        "UserName": usersett.userName(),
        "FullName": usersett.fullName(),
        "Email": usersett.email(),
        "Enable": statusBool,
        "Password": usersett.password(),
        "Role": dropRole.text(),
        "Potition":usersett.potition(),
        "LocationID": locationId,
        "LocationName": location.LocationName,
    }
    var url = "/usersetting/savedata";
    var validator = $("#AddUserSetting").data("kendoValidator");
    if (validator == undefined) {
        validator = $("#AddUserSetting").kendoValidator().data("kendoValidator");
    }
    if (validator.validate()) {
        // if (validator!=undefined) {  
        console.log("tes", validator)
        ajaxPost(url, param, function(res) {
            if (res.IsError != true) {
                $("#userModal").modal("hide");
                $("#nav-dex").css('z-index', 'none');
                usersett.ClearField();
                usersett.Reset();
                swal({
                    title:"Success!",
                    text: res.Message,
                    type: "success",
                    confirmButtonColor:"#3da09a"});
            } else {
                return swal({
                    title:"Error!",
                    text: res.Message,
                    type: "error",
                    confirmButtonColor:"#3da09a"});
            }
        });
    }
}

usersett.UpdateData = function() {
    var statusBool = $('#Status').bootstrapSwitch('state');
    var dropRole = $("#role").data("kendoDropDownList");
    var locationId = parseInt(msLocation.value())
    var location = _.find(msLocation.data(), function(e){return e.LocationID == locationId})
    // var param = {
    //     "Id" : usersett.get("Id"),
    //     "UserName": usersett.get("userName"),
    //     "FullName": usersett.get("fullName"),
    //     "Email": usersett.get("email"),
    //     "Enable": statusBool,
    //     "Password": usersett.get("password"),
    //     "Role": dropRole.text(),
    // }
    var param = {
        "Id": usersett.Id(),
        "UserName": usersett.userName(),
        "FullName": usersett.fullName(),
        "Email": usersett.email(),
        "Enable": statusBool,
        "Password": usersett.password(),
        "Role": dropRole.text(),
        "Potition":usersett.potition(),
        "LocationID": locationId,
        "LocationName": location.LocationName,
    }
    var url = "/usersetting/updatedata";
    var validator = $("#AddUserSetting").data("kendoValidator");
    if (validator == undefined) {
        validator = $("#AddUserSetting").kendoValidator().data("kendoValidator");
    }
    if (validator.validate()) {
        ajaxPost(url, param, function(res) {
            if (res.IsError != true) {
                // usersett.set("edit", false)
                usersett.edit(false)
                $("#userModal").modal("hide");
                $("#nav-dex").css('z-index', 'none');
                usersett.ClearField();
                usersett.Reset();
                swal({
                    title:"Success!",
                    text: res.Message,
                    type: "success",
                    confirmButtonColor:"#3da09a"});
            } else {
                return swal({
                    title:"Error!",
                    text: res.Message,
                    type: "error",
                    confirmButtonColor:"#3da09a"});
            }
        });
    }
}

usersett.EditData = function(idUser) {
    ajaxPost("/usersetting/getdata", {
        "Id": idUser
    }, function(res) {
        if (res.IsError != true) {
            // usersett.set("edit", true);
            usersett.edit(true)
            usersett.titleModel("Update User");

            $("#nav-dex").css('z-index', '0');

            $("#userModal").modal({
                backdrop: 'static',
                keyboard: false
            });

            var dataUser = res.Data.Records[0];
            $("#userModal").modal("show");
            usersett.Id(dataUser.Id);
            usersett.userName(dataUser.Username);
            usersett.fullName(dataUser.Fullname);
            usersett.email(dataUser.Email);
            usersett.password(dataUser.Password);
            usersett.confirmPassword(dataUser.Password);
            msLocation.value(dataUser.LocationID);
            msRoleVM.role(dataUser.Roles);
            usersett.potition(dataUser.Potition);
            // usersett.set("Id", dataUser.Id);
            // usersett.set("userName", dataUser.Username);
            // usersett.set("fullName", dataUser.Fullname);
            // usersett.set("email", dataUser.Email);
            // usersett.set("password", dataUser.Password);
            // usersett.set("confirmPassword", dataUser.Password);
            // msRoleVM.set("role", dataUser.Roles);

            $('#Status').bootstrapSwitch('state', dataUser.Enable);
            $("#role").data("kendoDropDownList").text(dataUser.Roles);
        } else {
            return swal({
            title:"Error!",
            text: res.Message,
            type: "error",
            confirmButtonColor:"#3da09a"});
        }
    });
}

usersett.Cancel = function() {
    // usersett.set("edit", false)
    usersett.edit(false)

    $("#userModal").modal("hide");
    $("#nav-dex").css('z-index', 'none');
    setTimeout(function() {
        usersett.ClearField();
    }, 500)

}

usersett.GetDataUser = function() {
    console.log(msUsernameVM.filterUser())
    // usersett.set("loading", false);
    usersett.edit(false)
    var dataSource = [];
    $("#MasterGridUser").html("");
    $("#MasterGridUser").kendoGrid({
        dataSource: {
            transport: {
                read: {
                    url: "/usersetting/getdata",
                    data: {
                        // "UserName" : msUsernameVM.get("filterUser"),
                        // "Role" : msRoleVM.get("filterRole"),
                        "UserName": msUsernameVM.filterUser(),
                        "Role": msRoleVM.filterRole(),
                        "Status": $('#StatusFilter').bootstrapSwitch('state')
                    },
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
                    // usersett.set("loading", false);
                    usersett.loading(false);

                    if (data.Data.Count == 0) {
                        console.log(data.Data.Records)
                        return dataSource;
                    } else {
                        return data.Data.Records;
                        console.log(data.Data.Records)
                    }
                },
                total: "Data.Count",
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
        columns: [{
                field: "Username",
                title: "User Name",
                width: 150,
                headerAttributes: {
                    class: 'k-header header-bgcolor'
                },
                // template: "# if(usersett.get('Edit') != 'false'){ # <a class='grid-select' id='ls' href='javascript:usersett.EditData(\"#: Id #\")'>#: Username #</a># }else{ #<div>#: Username #</div># } #"
                template: "# if(usersett.edit() != 'false'){ # <a class='onclickusersetting' id='ls' href='javascript:usersett.EditData(\"#: Id #\")'>#: Username #</a># }else{ #<div>#: Username #</div># } #"

            },
            {
                field: "Fullname",
                title: "Full Name",
                headerAttributes: {
                    class: 'k-header header-bgcolor'
                },
                width: 100

            },
            {
                /*field: "Enable",*/
                title: "Status",
                headerAttributes: {
                    class: 'k-header header-bgcolor'
                },
                template: function (x) {
                   // return x.Enable
                   if (x.Enable) {
                    return "Active"
                   }else{
                    return "InActive"
                   }
                },
                width: 50

            },
            {
                field: "Email",
                title: "Email",
                headerAttributes: {
                    class: 'k-header header-bgcolor'
                },
                width: 100

            },
            {
                field: "Roles",
                title: "Roles",
                headerAttributes: {
                    class: 'k-header header-bgcolor'
                },
                width: 100

            },
            {
                title: 'ACTION',
                width: 40,
                template: "<a href=\"javascript:usersett.delete('#: Id #', '#: Fullname #')\" class=\"btn btn-xs btn-danger\"><i class='fa fa-trash'></i></a>",
                attributes: {
                    class: 'align-center'
                }
            }
        ]
    });
}

usersett.onChangeUser = function(value) {
    msUsernameVM.filterUser(value)
    usersett.GetDataUser()
}

usersett.onChangerole = function(value) {
    msRoleVM.filterRole(value)
    usersett.GetDataUser()
}

usersett.getUserName = function() {
    // msUsernameVM.set("listUserName", []);
    msUsernameVM.listUserName([]);

    ajaxPost("/datamaster/getusername", {}, function(res) {
        var dataUser = Enumerable.From(res).OrderBy("$.username").ToArray();
        for (var u in dataUser) {
            msUsernameVM.listUserName.push({
                "text": dataUser[u].username,
                "value": dataUser[u].username,
            });
        }
    });
}

usersett.getRole = function() {
    var url = "/datamaster/getrolesrestricted";

    // msRoleVM.set("listRole", []);
    msRoleVM.listRole([]);

    ajaxPost("/datamaster/getrolesrestricted", {}, function(res) {
        console.log("getroles", res);
        var dataRole = Enumerable.From(res).OrderBy("$.name").ToArray();
        for (var r in dataRole) {
            msRoleVM.listRole.push({
                "text": dataRole[r].name,
                "value": dataRole[r].name,
            });
        }
    });
}

usersett.toggleFilter = function() {
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

    $('.k-grid').each(function(i, d) {
        try {
            $(d).data('kendoGrid').refresh();
        } catch (err) {}
    });

    $('.k-pivot').each(function(i, d) {
        $(d).data('kendoPivotGrid').refresh();
    });

    $('.k-chart').each(function(i, d) {
        $(d).data('kendoChart').redraw();
    });
    usersett.panel_relocated();
    // var FilterTitle = usersett.get("TitelFilter");
    var FilterTitle = usersett.TitelFilter();
    if (FilterTitle == " Hide Filter") {
        // usersett.set("TitelFilter", " Show Filter");
        usersett.TitelFilter(" Show Filter");
    } else {
        // usersett.TitelFilter(" Hide Filter");
        usersett.TitelFilter(" Hide Filter");
    }
}

usersett.panel_relocated = function() {
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

usersett.delete = function(id, name) {
console.log(id+"==="+name)

    model.Processing(true);
    swal({
        title: "Are you sure to delete " + name + "?",
        text: "Your will not be able to recover this user",
        type: "warning",
        showCancelButton: true,
        confirmButtonColor: "#3da09a",
        confirmButtonText: "Yes, delete it!",
        closeOnConfirm: false
    },
    function(res){
        if(res) {
            var url = "/datamaster/delete";
            var param = { id: id };
            ajaxPost(url, param, function(data){
                if(data != ""){
                    swal({
                        title:'Warning',
                        text: data,
                        type: 'error',
                        confirmButtonColor:"3da09a"});
                    model.Processing(false);
                }else{
                    swal({
                        title:'Success',
                        text: 'User has been deleted!',
                        type: 'success',
                        confirmButtonColor:"#3da09a"});
                    // customer.init()
                    usersett.GetDataUser();
                    model.Processing(false);
                }

            }, undefined);
        }
        else {
            model.Processing(false);
        }
    });
}
usersett.getDataLocation = function(){
    model.Processing(true)
    ajaxPost("/master/getdatalocation", {},function(res){
        model.Processing(false)   
        if (res.IsError){
            swal({
                title: "Error!",
                text: res.Message,
                type: "error",
                confirmButtonColor: "#3da09a"
            })
            return
        }
        _.each(res.Data, function(e){
            e.Text = e.LocationID +"-"+ e.LocationName
        })
        msLocation.data(res.Data)
    })
}
$(document).ready(function() {
    // kendo.bind($("#usersett"), usersett);
    // kendo.bind($("#select-username"), msUsernameVM)
    // kendo.bind($("#select-role"), msRoleVM)
    // kendo.bind($("#role"), msRoleVM)

    $('#StatusFilter').bootstrapSwitch('state', true);

    usersett.getUserName();
    usersett.getRole();
    usersett.GetDataUser();
    usersett.getDateNow();
    usersett.getDataLocation()
    // $("#select-username").kendoMultiSelect({
    //     valuePrimitive: true,
    //     filter: 'startswith', 
    //     dataTextField: 'text', 
    //     dataValueField: 'value'
    // })

    // $("#select-role").kendoMultiSelect({
    //     valuePrimitive: true,
    //     filter: 'startswith', 
    //     dataTextField: 'text', 
    //     dataValueField: 'value'
    // })

    // $("#role").kendoDropDownList({
    //     valuePrimitive: true,
    //     dataTextField: 'text',
    //     dataValueField: 'value',
    //     optionLabel: 'select roles'
    // })
    $(".k-icon.k-clear-value.k-i-close").hide()
});