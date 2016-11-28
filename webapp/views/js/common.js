
function setPageTitle(t){
	$("#pagetitle").text(t)
}

function GridConfig(){
	this.cfgObj =  {
            dataSource:{
                pageSize:10, 
                data: []
            },
            sortable:true,
			resizable:true,
            filterable:true,
			metadataUrl: "", 
            columns:[],
			error: ""
        };

	this.metadataUrl = "";
}

GridConfig.prototype.fetch = function(){
	thisObj = this;
	columns = [];

	if(thisObj.metadataUrl!=""){
		$.ajax({
			url:thisObj.metadataUrl,
			async:false
		}).done(function(data){
			if(data && data.Status=="OK"){
				columns = _.chain(data.Data.Fields).
					filter(function(item){
						return item.GridShow==true;
					}).
					sortBy(function(item){
						return item.GridColumn;
					}).
					map(function(item){
						return item.GridUseTemplate ?  {
							field:item.DBFieldID,
							title:item.Label,
							width:item.GridWidth,
							align:item.GridAlign,
							columnIndex:item.GridColumn
						} : {
							field:item.DBFieldID,
							title:item.Label,
							width:item.GridWidth,
							template:item.GridUseTemplate,
							columnIndex:item.GridColumn
						};
					}).
					value();
					//alert(kendo.stringify(columns));
			} else {
				alert("Error calling " + thisObj.metadataUrl);
			}
		}).fail(function(txt){
			alert("Error calling " + thisObj.metadataUrl);	
		});
		this.cfgObj.columns = columns;
	}
	return this.cfgObj;
}

GridConfig.prototype.set = function (attribute, value){
	this.cfgObj[attribute]=value;
	return this;
}

GridConfig.prototype.get = function (attribute){
	return this.cfgObj[attribute];
}

GridConfig.prototype.setDataSource = function (attribute, ds) {
	if(attribute==""){
		this.cfgObj.dataSource = ds;
	} else {
		this.cfgObj[attribute]=ds
	}
	return this;
}

/*
GridConfig.prototype.fromMetaData = function(mdts){
	var columns = [];
	mdts.forEach(function(obj,idx){
		columns.push({
			field:obj.DbFieldID,
			title:obj.Label
		})
	});
	this.set("columns",columns)
	return this;
}
*/

GridConfig.prototype.metadataFromUrl = function(url){
	//var models = metadataFromUrl(url);
	this.metadataUrl = url;
	return this;
}

function handleState(x, event){
    var target = event.target;
    var bindStateId = $(target).attr("data-bind");
    var stateObj = {};
    stateObj[bindStateId]=target.value;
    x.setState(stateObj);
}

function DataSource(){
	this.uri = "";
	this.data = [];
	this.uriType = "";
	this.postData = "";
	this.postDataType = "json";

	this.populate = function(data){

	};
}

DataSource.prototype.run = function(){
	thisObj = this;
	$.ajax({
		url: this.uri
	}).done(function(data){
		thisObj.populate(data);
	}).fail(function(txt){
		alert("Error call " + thisObj.metadataUrl);	
	});
}
