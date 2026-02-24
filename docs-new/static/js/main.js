/*clipboard*/

var getcodeelement = $('.clipboardjs'); /*create custom id*/

getcodeelement.each(function(i) {  
    /*target*/
    var currentId = 'codeblock' + (i + 1);
    $(this).attr('id', currentId);

    /*trigger*/
    var text = $(this).text();
    text = text.replace(/\$ /gi, '')
    .replace(/&/g, "&amp;")
    .replace(/</g, "&lt;")
    .replace(/>/g, "&gt;")
    .replace(/'/g, "&#39;")
    .replace(/"/g, "&quot;");
    var clipButton = '<div class="btn-copy-wrap"><button class="clipbtn" data-clipboard-text="' + text + '" data-clipboard-target="#' + currentId + '"><i class="far fa-copy"></i></button></div>';
       $(this).after(clipButton);
});

var clipboard = new Clipboard('.clipbtn');

/* Change copy icon to text when successfully copied*/
clipboard.on("success", (e)=>{
    console.info(e.trigger);
    console.info(e.trigger.childNodes[0]);
    let originalIcon = e.trigger.childNodes[0];

    var icon = e.trigger.childNodes[0];
    var text = document.createElement('span');
    text.textContent = "Copied!";
    text.style.color = "white";

    e.trigger.replaceChild(text, icon);

    setTimeout(()=>{
        e.trigger.replaceChild(originalIcon, text);
    },2000)
})
