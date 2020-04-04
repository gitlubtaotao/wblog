$(document).ready(function () {
    $('#pageSave').click(function (event) {
        event.preventDefault();
        $("#demo").text(simplemde.value());
        var tags = new Array();
        $(".tagButton .tagId").each(function (index, element) {
            tags.push($(element).text());
        });
        $("#tags").val(tags.join(","));
        var form = $("#postForm");
        $.ajax({
            url: form.attr('action'),
            type: form.attr('method'),
            dataType: "json",
            data: form.serializeArray(),
            success: function (data) {
                if (data['succeed']) {
                    toastr.success(data.message)
                } else {
                    toastr.error(data.message)
                }
            }
        });
    });
    //选择标签
    $('#select_tag').editable({
        mode: "inline",
        type: "select",
        select2: {
            placeholder: 'Select an option'
        },
        pk: 1,
        source: function () {
            var result = [];
            $.ajax({
                url: '/admin/tag/json',
                type: 'get',
                async: false,
                dataType: "json",
                success: function (data) {
                    result = data['data']
                }
            });
            return result
        },
        title: "select tag",
        display: function (value, response) {
            if (value !== null && typeof (response) !== "undefined") {
                var text = "";
                $.each(response, function (k, v) {
                    if (v['value'] === parseInt(value)) {
                        text = v['text'];
                        return false;
                    }
                });
                createButton(value, text);
            }
        }
    });
    $('#addTag').editable({
        mode: "inline",
        type: "text",
        pk: 1,
        url: "/admin/tag",
        placeholder: "add a tag",
        success: function (tag) {
            if (tag.succeed) {
                createButton(tag.data.ID, tag.data.name);
            } else {
                toastr.error(tag.message);
            }
        },
        error: function (e) {
            console.log(e);
        },
        display: function (value, response) {
            return false;   //disable this method
        }
    });
    $('#switchbtn').bootstrapSwitch({
        onText: '公开',
        offText: '不公开',
    });


    function createButton(tagId, tagName) {
        var button = `<button class="btn btn-default btn-sm tagButton">
                    <a href="/tag/` + tagId + `">` + tagName + `</a>
                    <a class="removeArticleTag " href="#" onclick="deleteTag(this);">
                        <span class="glyphicon glyphicon glyphicon-trash"></span>
                    </a>
                    <span class="tagId" hidden="hidden">` + tagId + `</span>
                    </button>&nbsp;`;

        $("#tag_content").append(button);
    }
});

function deleteTag(element) {
    $(element).parent().remove();
    // var id = $(element).next().text();
    // $.ajax({
    //     url: '/admin/tag/' + id,
    //     type: 'delete',
    //     dataType: 'json',
    //     success: function (data) {
    //         console.log(data)
    //     }
    // })
}



