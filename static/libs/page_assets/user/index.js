$(document).ready(function () {
    $('#example2').DataTable({
        'paging': true,
        'lengthChange': false,
        'searching': false,
        'ordering': true,
        'info': true,
        'autoWidth': false
    });
    $('.lock-button').on('click', function (e) {
        $.get($(e.target).data("href"), {}, function (data) {
            if (data.succeed) {
                toastr.success("Operation is successful")
            }else{
                toastr.error(data.message);
            }
        }, 'json');
    });


});