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
    $('tr').find('.published').on('click', function () {
        var id = $(this).data("id");
        var status = $(this).data("status");
        console.log(status);
        var _this = $(this);
        $.post("/admin/page/publish/" + id, {"_csrf": _csrf}, function (data) {
            console.log(data);
            if (data.succeed) {
                if (status === "false") {
                    _this.text('√')
                } else {
                    _this.text('x')
                }
                toastr.success("Operation is Successful")
            } else {
                toastr.error(data.message);
            }
        }, 'json')

    });


    $('#confirm-delete').on('show.bs.modal', function (e) {
        var _this = $(e.relatedTarget);
        $(this).find('.btn-ok').click(function () {
            $.ajax({
                url: _this.data('href')+"?_csrf="+_csrf,
                type: 'delete',
                data: {"_csrf": _csrf},
                dataType: 'json',
                success: function (data) {
                    if (data.succeed) {
                        $('#confirm-delete').modal('hide');
                        $('#tr_' + _this.data("id")).remove();
                        toastr.success('Operation is successful')
                    } else {
                        toastr.error(data.message);
                    }
                }
            });
        });
    });
});




