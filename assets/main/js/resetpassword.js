var model = {};
model.Mode = ko.observable("");
model.Processing = ko.observable(false);
var rp = {}

rp.resetpass = function(event) {
    event.preventDefault();
    var location = window.location.href;
    var tkn = location.substr(location.length - 30);
    if ($("#email-reset").val() != "") {
        if ($("#new_pass").val() != $("#conf_pass").val()) {
            $("#msg-error").html("Confirm New Password didn't match");
            return

        } else {
            var url = "/login/savenewpass";
            var param = {
                "Email": $("#email-reset").val(),
                "NewPass": hashPass($("#new_pass").val()),
                "UrlToLogin": window.location.origin,
                "Token": tkn,
            }
            console.log(window.location.origin)
            model.Mode("Process");
            rp.ShowHideLoader(true)
            ajaxPost(url, param, function(res) {
                if (res.Data.length === 0) {
                    $("#msg-error").html("Change Password Failed")
                    return swal("Error!", "Change Password Failed , please check your", "error");
                }
                rp.ShowHideLoader(false)
                swal("Success !", "Success Save New Password", "success");
                window.location.reload();
            });
            return false
        }

    } else {
        $("#msg-error").html("Please Fill The Email")
        return false
    }
}
rp.ShowHideLoader = function(show){
    if (show == true){
        $('#reset').hide();
        $('.loader').css('display','block');
    }else{
        $('#reset').show();
        $('.loader').css('display','none');
    }
}