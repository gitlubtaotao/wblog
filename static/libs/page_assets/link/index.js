$(document).ready(function () {
    $('#example2').DataTable({
        'paging': true,
        'lengthChange': false,
        'searching': false,
        'ordering': true,
        'info': true,
        'autoWidth': false
    });
    var dialog = $('#jquery_dialog');

    //新增link
    $('#add-dialog').on('click', function (e) {
        renderHtml("/admin/link")
    });

    $('.edit-row').on('click', function (e) {
        url = $(this).attr('data-href');
        renderHtml(url);
        var dataId = $(this).attr('data-id');
        $.get("/admin/link/" + dataId + '/show', {}, function (data) {
            if (data.succeed) {
                var form = dialog.find("#add-form");
                var link = data['link'];
                form.find('input[name="name"]').val(link['name']);
                form.find('input[name="url"]').val(link['url'])
            }
        }, "json")

    });

    $('#confirm-delete').on('show.bs.modal', function (e) {
        $(this).find('.btn-ok').click(function () {
            $.post($(e.relatedTarget).data('href'), {}, function (reuslt) {
                window.location.href = window.location.href;
            }, 'json');

        });
    });

    function renderHtml(url) {
        dialog.find(".modal-header").find("h3").text('Add/Edit Link');
        dialog.modal('show');
        dialog.find(".modal-body").empty().append($('#new_record').html());
        dialog.find('.btn-save').unbind("click"); //移除click
        dialog.find('.btn-save').click(function () {
            $.post(url,
                dialog.find('#add-form').serializeArray(), function (result) {
                    if (result.succeed) {
                        toastr.success("operation is successful");
                        window.location.reload();
                    } else {
                        toastr.error(result.message);
                    }
                }, 'json')
        });
    }

});