'use strict';

setInterval(() => {
    //debugger
    let c=document.getElementById("ghx-detail-head");
    if (c!==null){
        let divjiraibm=document.getElementById('jira_ibm');
        if (divjiraibm===null){
            let issuekey=document.getElementById('issuekey-val');
            let divJira=`<div id="jira_ibm">
            <a href="http://localhost:3000/app?peerreview=${issuekey.innerText}" target="_blank">PeerReview</a>
            </div>`;
            c.innerHTML+=divJira;
        }
    }
}, 1000);