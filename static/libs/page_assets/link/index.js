$(document).ready(function () {
    $('#example2').DataTable({
        'paging': true,
        'lengthChange': false,
        'searching': false,
        'ordering': true,
        'info': true,
        'autoWidth': false
    });

    $('#add-dialog').on('click', function (e) {
        var dialog = $('#jquery_dialog');
        dialog.modal('show');
        dialog.find(".modal-body").empty().append($('#new_record').html());
        var data = $('#example2').DataTable().row($(e.relatedTarget).parents('tr')).data();
        console.log(data);
        if (data) {
            var fields = $("#add-form").serializeArray();
            $.each(fields, function (i, field) {
                //jquery根据name属性查找
                $(":input[name='" + field.name + "']").val(data[i]);
            });
        }
        $(this).find('.btn-save').unbind("click"); //移除click
        $(this).find('.btn-save').click(function () {
            $.post($(e.relatedTarget).data('href'),
                $('#add-form').serialize(), function (result) {
                window.location.href = window.location.href;
            }, 'json')
        });
    });

    $('#confirm-delete').on('show.bs.modal', function (e) {
        $(this).find('.btn-ok').click(function () {
            $.post($(e.relatedTarget).data('href'), {}, function (reuslt) {
                window.location.href = window.location.href;
            }, 'json');

        });
    });
});