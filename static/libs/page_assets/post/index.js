$(document).ready(function () {
    var _csrf = $('.content').find('input[name="_csrf"]').val();
    $('#example2').DataTable({
        'paging': true,
        'lengthChange': false,
        'searching': false,
        'ordering': true,
        'info': true,
        'autoWidth': false
    });
    $('.publish-post').on('click', function () {
        var id = $(this).data('id');
        pushlish(id);
        function pushlish(id) {
            $.post("/admin/post/" + id + "/publish", {"_csrf": _csrf}, function (result) {
                console.log(result);
                if (result.succeed) {
                    toastr.success('Publish is successful')
                }
            }, "json")
        }
    });


    $('#confirm-delete').on('show.bs.modal', function (e) {
        $(this).find('.btn-ok').click(function () {
            $.post($(e.relatedTarget).data('href'), {"_csrf": _csrf}, function (result) {
                if (result['succeed']) {
                    $('tbody').find("#tr_" + result['id']).remove();
                    toastr.success("Delete post is successful");
                } else {
                    toastr.error(result['message']);
                }
                $('#confirm-delete').modal('hide')
            }, 'json');
        });
    });
});