/*!
 * jQuery Form Plugin 2.83 (11-JUL-2011)
 * Dual licensed under the MIT and GPL licenses:
 *   http://www.opensource.org/licenses/mit-license.php
 *   http://www.gnu.org/licenses/gpl.html
 */

/**
 * ajaxSubmit() provides a mechanism for immediately submitting
 * an HTML form using AJAX.
 */
$.fn.ajaxSubmit = function (options) {
    // fast fail if nothing selected
    if (!this.length) {
        return this;
    }

    var method, action, url, $form = this;

    if (typeof options == 'function') {
        options = {success: options};
    }

    method = this.attr('method');
    action = this.attr('action');
    url = (typeof action === 'string') ? $.trim(action) : '';
    url = url || window.location.href || '';
    if (url) {
        // clean url (don't include hash vaue)
        url = (url.match(/^([^#]+)/) || [])[1];
    }

    options = $.extend(true, {
        url: url,
        success: $.ajaxSettings.success,
        type: method || 'GET',
        iframeSrc: /^https/i.test(window.location.href || '') ? 'javascript:false' : 'about:blank'
    }, options);

    var n, v, a = this.formToArray(options.semantic);
    if (options.data) {
        options.extraData = options.data;
        for (n in options.data) {
            if (options.data[n] instanceof Array) {
                for (var k in options.data[n]) {
                    a.push({name: n, value: options.data[n][k]});
                }
            } else {
                v = options.data[n];
                v = $.isFunction(v) ? v() : v; // if value is fn, invoke it
                a.push({name: n, value: v});
            }
        }
    }

    var q = $.param(a);

    if (options.type.toUpperCase() == 'GET') {
        options.url += (options.url.indexOf('?') >= 0 ? '&' : '?') + q;
        options.data = null;  // data is null for 'get'
    } else {
        options.data = q; // data is the query string for 'post'
    }

    var callbacks = [];

    // perform a load on the target only if dataType is not provided
    if (!options.dataType && options.target) {
        var oldSuccess = options.success || function () {
        };
        callbacks.push(function (data) {
            var fn = options.replaceTarget ? 'replaceWith' : 'html';
            $(options.target)[fn](data).each(oldSuccess, arguments);
        });
    } else if (options.success) {
        callbacks.push(options.success);
    }

    options.success = function (data, status, xhr) { // jQuery 1.4+ passes xhr as 3rd arg
        var context = options.context || options;   // jQuery 1.4+ supports scope context
        for (var i = 0, max = callbacks.length; i < max; i++) {
            callbacks[i].apply(context, [data, status, xhr || $form, $form]);
        }
    };

    if ($.browser.msie && method == 'get') {
        var ieMeth = $form[0].getAttribute('method');
        if (typeof ieMeth === 'string')
            options.type = ieMeth;
    }
    $.ajax(options);

    // fire 'notify' event
    this.trigger('form-submit-notify', [this, options]);
    return this;
};

/**
 * ajaxForm() provides a mechanism for fully automating form submission.
 */
$.fn.ajaxForm = function (options) {
    if (this.length === 0) {
        if (!$.isReady && this.selector) {
            $(function () { //DOM not ready, queue
                $(this.selector, this.context).ajaxForm(options);
            });
        }
        return this;
    }

    return this.unbind('submit.form-plugin click.form-plugin')
        .bind('submit.form-plugin', function (e) {
            if (!e.isDefaultPrevented()) { // if event has been cancelled, don't proceed
                e.preventDefault();
                $(this).ajaxSubmit(options);
            }
        }).bind('click.form-plugin', function (e) {
            var target = e.target;
            var $el = $(target);
            if (!($el.is(":submit"))) {
                // is this a child element of the submit el?  (ex: a span within a button)
                var t = $el.closest(':submit');
                if (t.length == 0) {
                    return;
                }
                target = t[0];
            }
        });
};

/**
 * formToArray() gathers form element data into an array of objects that can
 * be passed to any of the following ajax functions: $.get, $.post, or load.
 * Each object in the array has both a 'name' and 'value' property.  An example of
 * an array for a simple login form might be:
 *
 * [ { name: 'username', value: 'jresig' }, { name: 'password', value: 'secret' } ]
 *
 * It is this array that is passed to pre-submit callback functions provided to the
 * ajaxSubmit() and ajaxForm() methods.
 */
$.fn.formToArray = function (semantic) {
    var a = [];
    if (this.length === 0) {
        return a;
    }

    var form = this[0];
    var els = semantic ? form.getElementsByTagName('*') : form.elements;
    if (!els) {
        return a;
    }

    var i, j, n, v, el, max, jmax;
    for (i = 0, max = els.length; i < max; i++) {
        el = els[i];
        n = el.name;
        if (!n) {
            continue;
        }

        v = $.fieldValue(el);
        if (v && v.constructor == Array) {
            for (j = 0, jmax = v.length; j < jmax; j++) {
                a.push({name: n, value: v[j]});
            }
        } else if (v !== null && typeof v != 'undefined') {
            a.push({name: n, value: v});
        }
    }
    return a;
};

/**
 * Returns the value of the field element.
 */
$.fieldValue = function (el) {
    var n = el.name, t = el.type, tag = el.tagName.toLowerCase();
    var successful = true;

    if (successful && (!n || el.disabled || t == 'reset' || t == 'button' ||
        (t == 'checkbox' || t == 'radio') && !el.checked ||
        (t == 'submit' || t == 'image') && el.form && el.form.clk != el ||
        tag == 'select' && el.selectedIndex == -1)) {
        return null;
    }

    if (tag == 'select') {
        var index = el.selectedIndex;
        if (index < 0) {
            return null;
        }
        var a = [], ops = el.options;
        var one = (t == 'select-one');
        var max = (one ? index + 1 : ops.length);
        for (var i = (one ? index : 0); i < max; i++) {
            var op = ops[i];
            if (op.selected) {
                var v = op.value;
                if (!v) { // extra pain for IE...
                    v = (op.attributes && op.attributes['value'] && !(op.attributes['value'].specified)) ? op.text : op.value;
                }
                if (one) {
                    return v;
                }
                a.push(v);
            }
        }
        return a;
    }
    return $(el).val();
}