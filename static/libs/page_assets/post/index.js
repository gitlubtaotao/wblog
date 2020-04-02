$(document).ready(function () {
    $('#example2').DataTable({
        'paging': true,
        'lengthChange': false,
        'searching': false,
        'ordering': true,
        'info': true,
        'autoWidth': false
    });
    function pushlish(id) {
        $.post("/admin/post/" + id + "/publish", {}, function (result) {
            console.log(result);
            if (result.succeed) {
                location.href = window.location.href;
            }
        }, "json")
    }

    $('#confirm-delete').on('show.bs.modal', function (e) {
        $(this).find('.btn-ok').click(function () {
            $.post($(e.relatedTarget).data('href'), {}, function (result) {

                location.href = window.location.href;
            }, 'json');

        });
    });
});