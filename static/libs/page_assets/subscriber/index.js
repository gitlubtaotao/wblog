$(document).ready(function () {
    $('#example2').DataTable({
        'paging': true,
        'lengthChange': false,
        'searching': false,
        'ordering': true,
        'info': true,
        'autoWidth': false
    });


    $('#confirm-delete').on('show.bs.modal', function (e) {
        $(this).find('.btn-ok').unbind("click");
        $(this).find('.btn-ok').click(function () {
            var subject = $('input[name="subject"]').val();
            var content = $('textarea[name="content"]').val();
            if (!subject) {
                alert("请填写主题");
                return;
            }
            if (!content) {
                alert("请填写内容");
                return;
            }

            $.post($(e.relatedTarget).data('href'), {subject: subject, content: content}, function (result) {
                console.log(result);
                if (result.succeed) {
                    alert("发送成功");
                } else {
                    alert(result.msg);
                }
                //window.location.href = window.location.href;
            }, 'json');
        });
    });
});