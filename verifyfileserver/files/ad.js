function animate(ele,target,spd){
 clearInterval(ele.timer);
 var speed = target>ele.offsetLeft?10:-10;
 ele.timer = setInterval(function () {
	 var val = target - ele.offsetLeft;
	 
	 if(Math.abs(val)<Math.abs(speed)){
	   ele.style.left = target + "px";
	   clearInterval(ele.timer);
	 }else{
		ele.style.left = ele.offsetLeft + speed + "px";
	 }
    },spd)
}
window.onload = function () {
  var all = document.getElementById("all");
  var screen = all.firstElementChild || all.firstChild;
  var imgWidth = screen.offsetWidth;
  var ul = screen.firstElementChild || screen.firstChild;
  var ol = screen.children[1];
  var div = screen.lastElementChild || screen.lastChild;
  var spanArr = div.children;
  var spd=10;
  var interval=2000;
  var ulNewLi = ul.children[0].cloneNode(true);
  ul.appendChild(ulNewLi);
  for(var i=0;i<ul.children.length-1;i++){
      var olNewLi = document.createElement("li");
      olNewLi.innerHTML = i+1;
      ol.appendChild(olNewLi)
  }
  var olLiArr = ol.children;
  olLiArr[0].className = "current";
  for(var i=0;i<olLiArr.length;i++){
      olLiArr[i].index = i;
      olLiArr[i].onmouseover = function () {
          for(var j=0;j<olLiArr.length;j++){
              olLiArr[j].className = "";
          }
          this.className = "current";
          key = square = this.index;
          animate(ul,-this.index*imgWidth,spd);
      }
  }
  var timer = setInterval(autoPlay,interval);
  var key = 0;
  var square = 0;
  function autoPlay(){
      key++;
      if(key>olLiArr.length){
          ul.style.left = 0;
          key = 1;
      }
      animate(ul,-key*imgWidth,spd);
      square++;
      if(square>olLiArr.length-1){
          square = 0;
      }
      for(var i=0;i<olLiArr.length;i++){
          olLiArr[i].className = "";
      }
      olLiArr[square].className = "current";
  }
  all.onmouseover = function () {
      div.style.display = "block";
      clearInterval(timer);
  }
  all.onmouseout = function () {
      div.style.display = "none";
      timer = setInterval(autoPlay,interval);
  }
  spanArr[0].onclick = function () {
      key--;
      if(key<0){
          ul.style.left = -imgWidth*(olLiArr.length)+"px";
          key = olLiArr.length-1;
      }
      animate(ul,-key*imgWidth,spd);
      square--;
      if(square<0){
          square = olLiArr.length-1;
      }
      for(var i=0;i<olLiArr.length;i++){
          olLiArr[i].className = "";
      }
      olLiArr[square].className = "current";
  }
  spanArr[1].onclick = function () {
      autoPlay();
  }
}