/*clipboard*/

const codeBlockElements = $('.clipboardjs'); /*create custom id*/

codeBlockElements.each(function(i) {  
    /*target*/
    var currentId = 'codeblock' + (i + 1);
    $(this).attr('id', currentId);

    /*trigger*/
    var text = $(this).text();
    text = text.replace(/\$ /gi, '').replace(/"/g, "&quot;");
    var clipButton = '<div class="btn-copy-wrap"><button class="clipbtn" data-clipboard-text="' + text + '" data-clipboard-target="#' + currentId + '"><i class="far fa-copy"></i></button></div>';
       $(this).after(clipButton);
});

var clipboard = new Clipboard('.clipbtn');

/* Change copy icon to text when successfully copied*/
clipboard.on("success", (e) => {
    const originalIcon = e.trigger.querySelector('i');
    if (!originalIcon) return;

    const feedbackText = document.createElement('span');
    feedbackText.textContent = "Copied!";
    feedbackText.classList.add('clipboard-success');

    e.trigger.replaceChild(feedbackText, originalIcon);

    setTimeout(() => {
        e.trigger.replaceChild(originalIcon, feedbackText);
    }, 2000);
});
