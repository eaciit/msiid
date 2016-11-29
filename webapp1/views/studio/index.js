setPageTitle("Studio")

var Button = React.createClass({
    render: function(){
        var btnClassName, btnGlyp;

        btnClassName = "btn btn-xs btn-" + this.props.type;
        btnGlyp = "glyphicon glyphicon-" + this.props.icon;
        return <button className={btnClassName} title={this.props.title}>
                    <span className={btnGlyp} />
                    &nbsp;
                    {this.props.title}
                </button>;
    }
})

var WidgetLists = React.createClass({
    render:function(){
        return <div className="widgetlists" _style={{padding:3+"px",backgroundColor:"whitesmoke"}}>
                    <div className="btn-toolbar">
                        <div className="btn-group">
                            <Button type="primary" icon="plus" title="Add Datasource" />
                        </div>
                        <div className="btn-group">
                            <Button type="primary" icon="th-large" title="Container" />
                        </div>
                        <div className="btn-group">
                            <Button type="primary" icon="th" title="Grid" />
                            <Button type="primary" icon="list-alt" title="Form" />
                            <Button type="primary" icon="search" title="Selector" />
                            <Button type="primary" icon="signal" title="Chart" />
                        </div>
                        <div className="btn-group">
                            <Button type="primary" icon="font" title="Label" />
                            <Button type="primary" icon="map-marker" title="Map" />
                            <Button type="primary" icon="picture" title="Multimedia" />
                            <Button type="primary" icon="pencil" title="Input" />
                            <Button type="primary" icon="inbox" title="Button" />
                        </div>
                        <div className="btn-group">
                            <Button type="primary" icon="fire" title="Script" />
                            <Button type="primary" icon="gift" title="Orchestration" />
                        </div>
                    </div>
            </div>;
    }
});

function addNode(data, idfieldname, parentid, childfieldname,obj){
    var node = getNode(data,idfieldname,parentid,childfieldname);
    if(node!=null){
        var widgets = node.widgets;
        if(!widgets || widgets.length==0){
            widgets=[];
        }
        obj.text = obj.title + " ["+obj.id+"]";
        widgets.push(obj);
        node.widgets = widgets;
    }
}

function getNode(data, fieldname, id, childnodename){
    var found = false;
    var l = data.length;

    for(var i=0;i<l;i++){
        var dataitem = data[i];
        if(dataitem[fieldname]==id){
            return dataitem;
        }

        if(!dataitem[childnodename]){
            var childnodes = dataitem[childnodename];
            var childnodefind = getNode(childnodes, fieldname, id, childnodename);
            if(childnodefind!=null){
                return childnodefind;
            }
        }
    }

    return null;
}

var Widget = React.createClass({
    getInitialState: function(){
        return {id:"",title:""};
    },
    render:function(){
        return <div className="widgetbox">
            </div>;
    }
})

var PageLayout = React.createClass({
    getInitialState: function(){
        return {
            pagestudio:{},
            dswidgetlists:[
                {id:"DS",text:"Data Source [DS]",widgets:[]},
                {id:"WI",text:"Widgets [WI]",widgets:[]}    
            ]};
    },
    componentDidMount: function(){
        var thisObj = $(this.refs.layoutwidgetlist);
        thisObj.kendoTreeView({
            //dataTextField: ["text"],
            select: function(e){
                console.log($(e.node).find(".k-in").text());
            }
        });
        /*
        this.initWidgetDataSource();
        this.addDataSource({id:"DS01",title:"HDFS"});
        this.addWidget({id:"WI01",title:"Container"});
        */
    },
    setPageStudio: function(s){
        //console.log("setting pagestudio state of layout object");
        this.setState({pagestudio:{widgets:100}});
    },
    addDataSource: function(ds){
        var dsdata = this.state.dswidgetlists;
        ds.text=ds.title + "["+ds.id+"]";
        addNode(dsdata,"id","DS","",ds);
        this.setState({dswidgetlists:dsdata});
        this.initWidgetDataSource();
    },
    addWidget: function(widget, parentid){
        var dsdata = this.state.dswidgetlists;
        if(!parentid || parentid==""){
            parentid="WI"
        }
        if(!widget.widgets){
            widget.widgets=[];
        }
        widget.text=widget.title + "["+widget.id+"]";
        addNode(dsdata,"id",parentid,"widgets",widget);
        this.setState({dswidgetlists:dsdata});

        this.initWidgetDataSource();
        this.refreshWidgetPreview();

        //console.log("Page studio:",this.state.pagestudio);
        if(this.state.pagestudio && this.state.pagestudio.state){
            var numofexistingwidget = this.state.pagestudio.state.widgets;
            this.state.pagestudio.setState({widgets:numofexistingwidget+1});
        }
    },
    refreshWidgetPreview: function(){
    },
    setWidgets:function(widgets){
        var thisObj = $(this.refs.layoutwidgetlist);
        thisObj.data("kendoTreeView").setDataSource(new kendo.data.HierarchicalDataSource({
                    data: widgets,
                    schema: {
                        model: {
                            children: "widgets"
                        }
                    }
                }));
    },
    initWidgetDataSource: function(){
        var thisObj = $(this.refs.layoutwidgetlist);
        thisObj.data("kendoTreeView").setDataSource(new kendo.data.HierarchicalDataSource({
                    data: this.state.dswidgetlists,
                    schema: {
                        model: {
                            children: "widgets"
                        }
                    }
                }));
    },
    render:function(){
        return <div className="pagelayout">
                <div ref="layoutwidgetlist" />
            </div>
    }
});

var PageStudio = React.createClass({
    getInitialState:function(){
        return {widgets:{
                count:1000
            }};
    },
    setWidgetCount:function(c){
        this.setState({widgets:c});
    },
    render:function(){
        var wh = ($(window).height()-100) + "px";
        return <div className="pagestudio" 
                //style={{border:"solid 1px #aaa",padding:"2px",height:wh}}
                style={{padding:"2px"}}
            >
            &nbsp;{this.state.widgets.count}&nbsp;
            </div>
    }
});

var StudioFormApp = React.createClass({
    getInitialState:function(){
        return {
            widgets: []
        };
    },
    refreshWidgetList:function(){
        this.refs.pagelayout.setWidgets(this.state.widgets);
    },
    addWidget:function(widget, parentid){
        this.refreshWidgetList();
    },
    removeWidget:function(widgetid){
        this.refreshWidgetList();
    },
    addDataSource:function(ds){
        this.refreshWidgetList();
    },
    removeDataSource:function(dsid){
        this.refreshWidgetList();
    },
    componentDidMount:function(){
        var layout = this.refs.pagelayout;
        var studio = this.refs.pagestudio;
        
        this.addDataSource({id:"DS01",title:"HDFS"});
        this.addWidget({id:"WI01",title:"Container"});
    },
    render:function(){
            return <div> 
                <div className="row" style={{marginBottom:"5px"}}>
                    <div className="col-md-12">
                        <WidgetLists ref="widgetlists" />
                    </div>
                </div>
                <div className="row">
                    <div className="col-md-2">
                        <PageLayout ref="pagelayout" />
                    </div>
                    <div className="col-md-10">
                        <PageStudio ref="pagestudio" />
                    </div>
                </div>
            </div>
    }
});

ReactDOM.render(<StudioFormApp />, document.getElementById("panel_form"));