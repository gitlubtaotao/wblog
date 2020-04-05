$(document).ready(function () {
    var _csrf = $('form').find('input[name="_csrf"]').val();
    $('#unbind_github').on("click", function () {

        $.post("/admin/profile/github/unbind", {"_csrf": _csrf}, function (result) {
            toastr.success(result.message);
            if (result.succeed) {
                window.location.reload(true)
            }
        }, 'json');
    });

    $('#unbind_email').on('click', function () {

    });
    $('#bind_email').on('click', function () {
        var email = $('#inputEmail').val();
        if (email === "") {
            toastr.error("email is not null");
            return;
        }
        $.post("/admin/profile/email/bind", {email: email,"_csrf": _csrf}, function (data) {
            if (data.succeed) {
                toastr.success(data.message);
            } else {
                toastr.error(data.message)
            }
        })
    });

    $('#submit_form').on('click', function () {
        var form = $('.form-horizontal');
        var password = form.find('input[name="password"]').val();
        console.log(password);
        var confirm_password = form.find('input[name="confirm_password"]').val();
        if (password !== "") {
            if (password !== confirm_password) {
                toastr.error("Password is error");
                return
            }
        }
        $.ajax({
            url: form.attr("action"),
            type: form.attr("method"),
            data: form.serializeArray(),
            dataType: 'json',
            success: function (data) {
                if (data.succeed) {
                    toastr.success("Save is successful");
                } else {
                    toastr.error(data.message);
                }
            }
        })
    });
    //上传图片
    $('#upload_avatar').on("change", function () {
        var formData = new FormData();
        formData.append('file', $(this)[0].files[0]);
        formData.append("_csrf",_csrf);
        $.ajax({
            url: '/admin/upload',
            type: 'POST',
            cache: false,
            data: formData,
            processData: false,
            contentType: false,
        }).done(function (res) {
            if (res.succeed) {
                $('input[name="avatar_url"]').val(res.url);
                $('input[name="secret_key"]').val(res.key);
                $('#show_avatar').attr("src", res.url);
            } else {
                toastr.error("Upload avatar url is error")
            }
        }).fail(function (res) {

        });
    });
});