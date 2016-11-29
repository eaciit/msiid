var EC = {}
window.EC = EC;

EC.DataBrowser = React.createClass({
    getInitialState: function(){
        return {
            simpleQuery:"", gridConfig:{}
        };
    },
    componentDidMount: function(){
        this.initGrid(this.props.gridConfig);
    },
    handleState : function(event){
        handleState(this,event);
    },
    hide: function(){
        $(this.refs.main).hide();
    },
    show: function(){
        $(this.refs.main).show();
    },
    populate: function(){

    },
    initGrid: function(cfg){
        var mainObj = $(this.refs.grid);
        if(mainObj.data("kendoGrid")!=undefined){
            mainObj.data("kendoGrid").destroy();
        }
        mainObj.kendoGrid(cfg);
    },
    setData: function(data, currentPage, pageCount){
        //this.refs.grid.setData(data, currentPage, pageCount)
    },
    refresh: function(){
        var q = this.state.simpleQuery;
        var refreshQuery = {
            q: q
        };
        alert("Querying: " + q);
    },
    add: function(){
        alert("This is to run add operation");
    },
    delete: function(){
        confirm("Are you sure you want to delete selected data ?");
    },

    manage: function(){
    },

    render: function(){
        var searchControls, buttonControls, divControls;

        if(this.props.hideSearch!="true" && this.props.searchBox!=undefined){
            searchControls = this.props.searchBox;
        } else {
            searchControls = this.props.hideSearch=="true" ? <div></div> :  
                            <div className="input-group input-group-sm">
                                <input type="text" className="form-control" ref="bind_simpleQuery" 
                                    placeholder="Search for ..." 
                                    data-bind="simpleQuery" value={this.state.simpleQuery} 
                                    onChange={this.handleState}
                                    />
                                <span className="input-group-btn">
                                    <button className="btn btn-sm btn-default" type="button" 
                                        onClick={this.refresh}>Refresh</button>
                                    <button className="btn btn-sm btn-default">
                                        <span className="glyphicon glyphicon-tasks"></span>
                                    </button>
                                </span>
                            </div>;
        }

        buttonControls =  this.props.hideButton=="true" ? <div></div> : 
                        <div>
                        <div className="btn-group">
                            <button type="button" className="btn btn-sm btn-primary" onClick={this.add}>Add New</button>
                            <button type="button" className="btn btn-sm btn-warning" onClick={this.delete}>Delete</button>
                        </div>
                        &nbsp;&nbsp;
                        <button type="button" className="btn btn-sm btn-primary" onClick={this.manage}>
                            <span className="glyphicon glyphicon-cog"></span>
                        </button>
                        </div>;

        if(this.props.searchBox==undefined){
            divControls = this.props.hideSearch=="true" ? 
                        <div style={{marginBottom:5+"px"}} className="row">
                        <div className="col-sm-12">
                            {buttonControls}
                        </div>
                        </div> :
                        <div style={{marginBottom:5+"px"}} className="row">
                        <div className="col-sm-6">
                            {searchControls}
                        </div>
                        <div className="col-sm-6">
                            {buttonControls}
                        </div>
                        </div>;
        }else{
            divControls = this.props.hideSearch=="true" ? 
                        <div style={{marginBottom:5+"px"}} className="row">
                            <div className="col-sm-12">
                                {buttonControls}
                            </div>
                        </div> :
                        <div>
                            <div style={{marginBottom:5+"px"}} className="row">
                                <div className="col-sm-12">
                                    {searchControls}
                                </div>
                            </div>
                            <div style={{marginBottom:5+"px"}} className="row">
                                <div className="col-sm-12">
                                    {buttonControls}
                                </div>
                            </div>
                        </div>;
        };

        return <div ref="main">
                {divControls}
                <div className="ecgrid" ref="grid"></div>
            </div>;
    }
});

