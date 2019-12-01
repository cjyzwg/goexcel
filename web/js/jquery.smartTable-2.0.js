 /*!
 * Collapsible, jQuery Plugin
 * jquery.smartTable-2.0.js
 *
 * Copyright (c) 2013
 * author: Bian Kaiming (Walle)
 * @ version: 1.0.1
 *
 * Date: Fri Nov 14 15:12:11 2013
 */
;(function($){
	$.fn.smartTable = function(options){
		var defaultVal = {
			editText:true, 																		//是否可编辑单元格内容
			editRow :true, 																		//是否可操作单元行，添加或删除
			addBtnPosition:"bottom",															//添加按钮的位置，left和bottom，默认bottom，left可动态从中间插入，bottom只能从下方添加
			addBtnBg:["url(addbtn.png) no-repeat","url(addbtn_hover.png)"],						//添加按钮自定义，若覆盖和移除背景相同，两项需写一样的内容
			addBtnWidth:"38px",																	//添加按钮自定义宽度
			addBtnHeight:"11px",																//添加按钮自定义高度
			delBtnBg:["url(delete.png) 2px 6px no-repeat", "url(delete.png) 2px 6px no-repeat"],//删除按钮自定义
			delBtnWidth:"20px",																	//删除按钮自定义宽度
			tableContent:[ 																		//表格内json数据,可为空
				{th:"表头1",td:["1-1", "1-2", "1-3" ,"1-4", "1-5"]},
				{th:"表头2",td:["2-1", "2-2", "2-3" ,"2-4", "2-5"]},
				{th:"表头3",td:["3-1", "3-2", "3-3" ,"3-4", "3-5"]},
				{th:"表头4",td:["4-1", "4-2", "4-3" ,"4-4", "4-5"]},
				{th:"表头5",td:["5-1", "5-2", "5-3" ,"5-4", "5-5"]},
				{th:"表头6",td:["6-1", "6-2", "6-3" ,"6-4", "6-5"]}
			],
			checkzerocontent:false
		};
		var options = $.extend(defaultVal,options);	
		
		return this.each(function(){
		
			var self = $(this);//取自身对象
			var divMark;//取自身对象标记
			if($(self).attr("id") != undefined){
				divMark = $(self).attr("id");
			}else if($(self).attr("class") != undefined){
				divMark = $(self).attr("class");
			};
			
			/* 初始化表格 */
			var init = function(){
				if(options.tableContent instanceof Array){
					$(self).append("<div class='tableDiv-"+divMark+"' style='position:static'><table><thead></thead><tbody></tbody></table></div>");//构建雏形
					//表头
					var oThead = "<tr>";
					var oTrMaxLength = new Number;
					for(var i = 0; i < options.tableContent.length; i++){
						oThead = oThead + "<th>"+options.tableContent[i].th+"</th>";
						//顺便找出最长的td
						if(oTrMaxLength < options.tableContent[i].td.length){
							oTrMaxLength = options.tableContent[i].td.length;
						};
					};
					oThead = oThead + "</tr>";
					//插入表头
					$(self).find("thead").append(oThead);
					//表内容
					var oTbody;
					for(var j = 0; j <= oTrMaxLength - 1; j++){
						oTbody = "<tr>";
						for(var i = 0; i < options.tableContent.length; i++){
							if(options.tableContent[i].td[j] != undefined){
								var tdstr = "<td>"+options.tableContent[i].td[j]+"</td>";
								if(options.tableContent[i].td[j]==0 && options.checkzerocontent==true){
									tdstr = "<td bgcolor='#FF0000'>"+options.tableContent[i].td[j]+"</td>";
								}
								oTbody = oTbody + tdstr;
								// oTbody = oTbody + "<td>"+options.tableContent[i].td[j]+"</td>";
							}else{
								var tdstr = "<td>--</td>";

								if(options.checkzerocontent==true){
									tdstr = "<td bgcolor='#FF0000'>--</td>";
								}
								oTbody = oTbody + tdstr

							};
						};
						oTbody = oTbody + "</tr>";
						$(self).find("tbody").append(oTbody);//分别插入表内容
					};
					
				};
			};
			init();//执行初始化

			
			/* 双击修改表格内容 */
			var editText = function(){
				$(self).find("td").live("dblclick",function(){
					var x = $(this).index();
					var y = $(this).parent().index();//alert("x:"+x+"--y:"+y);
					var tdText = $(this).text();
					var tdWidth = $(this).width();
					var tdHeight = $(this).height();
					$(this).html("<input type=\"text\" value=\""+tdText+"\" style=\"height:"+tdHeight+"px;line-height:"+tdHeight+"px;width:"+tdWidth+"px;\">");
					
					//获得焦点 并设置回车反应
					if($(this).find("input").focus()){
						$(document).keydown(function(e){
							if(e.keyCode == 13){
								$(self).find("input").parent().text($(self).find("input").val());
								btnPosition.addPosition();
							};
						});
					};
					//失去焦点
					$(this).find("input").blur(function(){
						$(this).parent().text($(this).val());
						btnPosition.addPosition();
					});
					btnPosition.addPosition();//双击后变化 添加 按钮位置
				});
			};
			if(options.editText == true){
				editText();//执行绑定双击事件
			}else{
				return false;
			};
			
			/* 动态添加行和删除行按钮位置 */
			var trIndex;//建立tr索引
			var editRow = function(){
				//添加按钮位置
				if(options.addBtnPosition == "left"){
					//再议
				}else{
					//添加按钮
					var addBtn = "<div id='addTr"+divMark+"' style='"+
								 "position:absolute;"+
								 "background:"+options.addBtnBg[0]+";"+
								 //"top:"+btnPosition.addBtnTop()+"px;"+
								 "left:"+btnPosition.addBtnLeft()+"px;"+
								 "width:"+options.addBtnWidth+";"+
								 "height:"+options.addBtnHeight+";"+
								 "display:none;"+
								 "'></div>";
					$(self).find(".tableDiv-"+divMark).append(addBtn);//add插入页面
					$("#addTr"+divMark).mouseover(function(){//覆盖显示
						$(this).css({"background":options.addBtnBg[1],"cursor":"pointer"});
					}).mouseleave(function(){
						$(this).css({"background":options.addBtnBg[0]});
					});
					
					//删除按钮
					var delBtn = "<div id='delTr"+divMark+"' style='display:none;position:absolute;width:"+options.delBtnWidth+";'><span style='"+
								 "display:block;"+
								 "float:right;"+
								 "background:"+options.delBtnBg[0]+";"+
								 "height:100%;"+
								 "width:"+options.delBtnWidth+"'></span></div>";
					$(self).find(".tableDiv-"+divMark).append(delBtn);//del插入页面
					//事件
					$(self).find("tbody tr").live("mouseover",function(){
						trIndex = $(this).index();//获得tr索引
						//按钮定位
						$("#delTr"+divMark).css({"left":$(this).offset().left+$(this).width()+"px","height":$(this).height()+"px","display":"block"});
						if($(self).css("position") == "static"){
							$("#delTr"+divMark).css({"top":$(this).offset().top +"px"});
						}else{
							$("#delTr"+divMark).css({"top":$(this).offset().top - $(self).offset().top +"px"});
						};
						$("#delTr"+divMark+">span").css({"background-position-y":$(this).height()/2-8+"px"});
						//判断是否最后一行
						if($(self).find("tbody tr").last().index() == 0){
							$("#delTr"+divMark).css({"display":"none"});
						};
					}).live("mouseleave",function(){
						$("#delTr"+divMark).css({"display":"none"});
					});

					//鼠标移动到按钮上
					$("#delTr"+divMark).mouseover(function(){
						$(this).css({"display":"block","cursor":"pointer"});
					}).mouseleave(function(){
						$(this).css({"display":"none"});
					});
					
				};
			};
			var btnPosition = {//表格和添加按钮位置对象
				tableTop : function(){return $(self).find("table").offset().top;},
				tableLeft : function(){return $(self).find("table").offset().left;},
				//addBtnTop : function(){return this.tableTop() + $(self).find("table").height()+3;},
				addBtnLeft : function(){return this.tableLeft() + ($(self).find("table").width()+2)/2-19;},
				addPosition : function(){$("#addTr"+divMark).css({"left":this.addBtnLeft()+"px"});}//"top":this.addBtnTop()+"px", 
			};
			if(options.editRow == true){
				editRow();//执行动态添加和删除行操作
				$(window).resize(function(){//窗口重载定位add按钮
					btnPosition.addPosition();
				});
			}else{
				return false;
			};
			
			/* 动态添加行执行 */
			$(self).find(".tableDiv-"+divMark).css("padding-bottom",options.addBtnHeight);
			$(self).find(".tableDiv-"+divMark).on("mouseover",function(){
				$("#addTr"+divMark).show();
				btnPosition.addPosition();
			}).on("mouseout",function(){
				$("#addTr"+divMark).hide();
			});//鼠标移动显示添加按钮
			$("#addTr"+divMark).on("click",function(){
				var addTd;
				for(var i = 0; i < options.tableContent.length; i++){
					addTd = addTd + "<td>--</td>";
				};
				$(self).find("table tbody").append("<tr>"+addTd+"</tr>");
				btnPosition.addPosition();
			});
			
			/* 动态删除行执行 */
			$("#delTr"+divMark).on("click",function(){
				$(self).find("tbody tr:eq("+trIndex+")").css("background-color","#fbd0d0").fadeOut(400,function(){
					$(this).remove();
					btnPosition.addPosition();//单击后变化 添加 按钮位置
					$("#delTr"+divMark).css({"display":"none"});//隐藏删除按钮
				});
			});
			var tableJSON;
			//获取数据函数
			var catchData = function(classname){
				var row = $("."+classname).find("table thead th");
				var col = $("."+classname).find("table tbody tr");
				tableJSON = '{"data":[';
				for(var i = 0; i <= row.index(); i++){
					tableJSON = tableJSON + '{"th":' + '"' + row.eq(i).text() + '"' +',"td":[';
					for(var j = 0; j <= col.index(); j++){
						if(j == col.index()){
							tableJSON = tableJSON + '"' +col.eq(j).find("td").eq(i).text() + '"';
						}else{
							tableJSON = tableJSON + '"'  + col.eq(j).find("td").eq(i).text() + '"' + ',';
						};
					};
					if(i == row.index()){
						tableJSON = tableJSON + ']}'
					}else{
						tableJSON = tableJSON + ']},';
					};
				};
				tableJSON = tableJSON + ']}'
			};
			/* 建立获取数据接口 */
			$.fn.smartTable.getData = function(classname){
				catchData(classname);
				return $.parseJSON(tableJSON);
				//return tableJSON;
			};
		
		});
	};
})(jQuery);


