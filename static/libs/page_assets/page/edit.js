$(document).ready(function () {
    $('#pageSave').click(function (event) {
        event.preventDefault();
        var form = $("#pageForm");
        var title = form.find('input[name="title"]').val();
        if (title === "") {
            toastr.error("Title is required");
            return false;
        }
        $("#demo").text(window.simpleMde.value());
        $.post(form.attr('action'), form.serializeArray(), function (data) {
            if (data.succeed) {
                toastr.success("Operation is successful");
            } else {
                toastr.error(data.message);
            }
        }, 'json');

    });
    $('#switchbtn').bootstrapSwitch({
        onText: '公开',
        offText: '不公开',
    });
});