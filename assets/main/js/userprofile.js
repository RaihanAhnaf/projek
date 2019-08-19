var profile = {}
var UserObj = {
    Id: '',
    Username: '',
    Fullname: '',
    CellularNo: '',
    Password: '',
    Email: '',
    Roles: '',
    LastLogin: '',
    Enable: true
};
profile.UserModel = ko.observable(UserObj)
profile.EditIsTrue = ko.observable(false)
profile.DatePageBar = ko.observable()

profile.getDateNow = function () {
    var page = moment(new Date()).format("DD MMMM YYYY HH:mm")
    profile.DatePageBar(page)
}

model.Processing(false)
profile.GetProfile = function () {
    var url = "/userprofile/getprofile"
    var param = {
        username: userinfo.usernameh()
    }
    ajaxPost(url, param, function (res) {
        var data = res.data[0]
        profile.UserModel({
            Id: data.Id,
            Username: data.Username,
            Fullname: data.Fullname,
            CellularNo: data.CellularNo,
            Password: data.Password,
            Email: data.Email,
            Roles: data.Roles,
            LastLogin: data.LastLogin,
            Enable: data.Enable
        })
    })

}
profile.EnableField = function () {
    $("#Name").prop("disabled", false);
    $("#username").prop("disabled", false);
    $("#celluler").prop("disabled", false);
    $("#email").prop("disabled", false);
    $("#password").prop("disabled", false);
    $("#enable").prop("disabled", false);

}
profile.EditProfile = function () {
    profile.EnableField()
    profile.EditIsTrue(true)

}
profile.EditCanceled = function () {
    model.Processing(true)
    setTimeout(function () {
        location.reload()
    }, 1000);

}
profile.SaveProfile = function () {
    var pass = $("#password").val()
    var repass = $("#repassword").val()
    if (pass != repass) {
        return swal({
            title:"Re-Password is not the same as Password!", 
            text:"Please check your password",
            type: "error",
            confirmButtonColor:"#3da09a"})
    } 
    var url = "/userprofile/saveprofile"
    var param = {
        Id: profile.UserModel().Id,
        Name: $("#Name").val(),
        Username: $("#username").val(),
        Celluler: $("#celluler").val(),
        Email: $("#email").val(),
        Password: pass
    }
    
    swal({
        title: "Are you sure?",
        text: "You want to save your profile!",
        type: "warning",
        showCancelButton: true,
        confirmButtonColor: "#3da09a",
        confirmButtonText: "Yes, do it!",
        cancelButtonText: "No, cancel !",
        closeOnConfirm: false,
        closeOnCancel: false
    }, function (isConfirm) {
        if (isConfirm) {
            model.Processing(true)
            ajaxPost(url, param, function (res) {
                console.log(res);
                if (res.IsError) {
                    swal({
                        title:"Failed",
                        text: res.Data != null ? res.Data : "Your profile is not saved",
                        type:  "error",
                        confirmButtonColor:"#3da09a",
                    });
                }
                else {
                    swal({
                        title:"Success!",
                        text: "Profile has been changed, you will need to login again.",
                        type: "success",
                        confirmButtonColor:"#3da09a"})
                    setTimeout(function () {
                        localStorage.clear();
                        window.location.href = "/logout/do";
                    }, 2000) 
                }       
            })
        } else {
            swal({
                title:"Cancelled",
                text: "Your profile is not saved",
                type:  "error",
                confirmButtonColor:"#3da09a",
            });
        }
        model.Processing(false)
    });
    

}
$(document).ready(function () {
    profile.GetProfile()
    profile.getDateNow()
});
