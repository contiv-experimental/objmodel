/******/ (function(modules) { // webpackBootstrap
/******/ 	// The module cache
/******/ 	var installedModules = {};

/******/ 	// The require function
/******/ 	function __webpack_require__(moduleId) {

/******/ 		// Check if module is in cache
/******/ 		if(installedModules[moduleId])
/******/ 			return installedModules[moduleId].exports;

/******/ 		// Create a new module (and put it into the cache)
/******/ 		var module = installedModules[moduleId] = {
/******/ 			exports: {},
/******/ 			id: moduleId,
/******/ 			loaded: false
/******/ 		};

/******/ 		// Execute the module function
/******/ 		modules[moduleId].call(module.exports, module, module.exports, __webpack_require__);

/******/ 		// Flag the module as loaded
/******/ 		module.loaded = true;

/******/ 		// Return the exports of the module
/******/ 		return module.exports;
/******/ 	}


/******/ 	// expose the modules object (__webpack_modules__)
/******/ 	__webpack_require__.m = modules;

/******/ 	// expose the module cache
/******/ 	__webpack_require__.c = installedModules;

/******/ 	// __webpack_public_path__
/******/ 	__webpack_require__.p = "";

/******/ 	// Load entry module and return exports
/******/ 	return __webpack_require__(0);
/******/ })
/************************************************************************/
/******/ ([
/* 0 */
/***/ function(module, exports, __webpack_require__) {

	/** @jsx React.DOM *//** @jsx React.DOM */

	// Little hack to make ReactBootstrap components visible globally
	Object.keys(ReactBootstrap).forEach(function (name) {
	    window[name] = ReactBootstrap[name];
	});

	// Navigation tab
	var ControlledTabArea = __webpack_require__(1)

	// Render the main tabs
	React.render(React.createElement(ControlledTabArea, null), document.getElementById('mainViewContainer'));


/***/ },
/* 1 */
/***/ function(module, exports, __webpack_require__) {

	/** @jsx React.DOM */// navTab.js
	// Navigation tab

	// panels
	var HomePane = __webpack_require__(2)
	var NetworkPane = __webpack_require__(3)
	var GroupsPane = __webpack_require__(4)
	var PolicyPane = __webpack_require__(5)
	var VolumesPane = __webpack_require__(6)

	// Define tabs
	var ControlledTabArea = React.createClass({displayName: "ControlledTabArea",
	  getInitialState: function() {
	    return {
	      key: 1,
	    };
	  },

	  getStateFromServer: function() {
	    // Sort function for all contiv objects
	    var sortObjFunc = function(first, second) {
	      if (first.key > second.key) {
	          return 1
	      } else if (first.key < second.key) {
	          return -1
	      }

	      return 0
	    }

	    // Get all endpoints
	    $.ajax({
	      url: "/endpoints",
	      dataType: 'json',
	      success: function(data) {

	        // Sort the data
	        data = data.sort(sortObjFunc);

	        this.setState({endpoints: data});

	        // Save it in a global variable for debug
	        window.globalEndpoints = data
	      }.bind(this),
	      error: function(xhr, status, err) {
	        console.error("/endpoints", status, err.toString());
	      }.bind(this)
	    });

	    // Get all networks
	    $.ajax({
	      url: "/api/networks/",
	      dataType: 'json',
	      success: function(data) {

	        // Sort the data
	        data = data.sort(sortObjFunc);

	        this.setState({networks: data});

	        // Save it in a global variable for debug
	        window.globalNetworks = data
	      }.bind(this),
	      error: function(xhr, status, err) {
	        console.error("/api/networks/", status, err.toString());
	      }.bind(this)
	    });

	    // Get all endpoint groups
	    $.ajax({
	      url: "/api/endpointGroups/",
	      dataType: 'json',
	      success: function(data) {

	        // Sort the data
	        data = data.sort(sortObjFunc);

	        this.setState({endpointGroups: data});

	        // Save it in a global variable for debug
	        window.globalEndpointGroups = data
	      }.bind(this),
	      error: function(xhr, status, err) {
	        console.error("/api/endpointGroups/", status, err.toString());
	      }.bind(this)
	    });

	    // Get all policies
	    $.ajax({
	      url: "/api/policys/",
	      dataType: 'json',
	      success: function(data) {

	        // Sort the data
	        data = data.sort(sortObjFunc);

	        this.setState({policies: data});

	        // Save it in a global variable for debug
	        window.globalPolicies = data
	      }.bind(this),
	      error: function(xhr, status, err) {
	        console.error("/api/policys/", status, err.toString());
	      }.bind(this)
	    });

	    // Get all rules
	    $.ajax({
	      url: "/api/rules/",
	      dataType: 'json',
	      success: function(data) {

	        // Sort the data
	        data = data.sort(sortObjFunc);

	        this.setState({rules: data});

	        // Save it in a global variable for debug
	        window.globalRules = data
	      }.bind(this),
	      error: function(xhr, status, err) {
	        console.error("/api/rules/", status, err.toString());
	      }.bind(this)
	    });
	  },
	  componentDidMount: function() {
	    this.getStateFromServer();

	    // Get state every 2 sec
	    setInterval(this.getStateFromServer, 2000);
	  },
	  handleSelect: function(key) {
	    console.log('selected Tab ' + key);
	    this.setState({key: key});
	  },

	  render: function() {
	      var self = this

	    return (
	      React.createElement(TabbedArea, {activeKey: this.state.key, onSelect: this.handleSelect}, 
	        React.createElement(TabPane, {eventKey: 1, tab: "Home"}, 
	            React.createElement(HomePane, {key: "home", endpoints: this.state.endpoints})
	        ), 
	        React.createElement(TabPane, {eventKey: 3, tab: "Networks"}, " ", React.createElement("h3", null, " Networks "), 
	            React.createElement(NetworkPane, {key: "networks", networks: this.state.networks})
	        ), 
	        React.createElement(TabPane, {eventKey: 4, tab: "Groups"}, " ", React.createElement("h3", null, " Groups "), 
	            React.createElement(GroupsPane, {key: "groups", endpointGroups: this.state.endpointGroups})
	        ), 
	        React.createElement(TabPane, {eventKey: 5, tab: "Policies"}, " ", React.createElement("h3", null, " Policy "), 
	            React.createElement(PolicyPane, {key: "policy", policies: this.state.policies})
	        ), 
	        React.createElement(TabPane, {eventKey: 6, tab: "Volumes"}, " ", React.createElement("h3", null, " Volumes "), 
	            React.createElement(VolumesPane, {key: "volumes", volumes: this.state.volumes})
	        )
	      )
	    );
	  }
	});

	module.exports = ControlledTabArea


/***/ },
/* 2 */
/***/ function(module, exports, __webpack_require__) {

	/** @jsx React.DOM */// home.js
	// Display Endpoint information

	var HomePane = React.createClass({displayName: "HomePane",
	  	render: function() {
			var self = this

			if (self.props.endpoints === undefined) {
				return React.createElement("div", null, " ")
			}

			// Walk thru all the endpoints
			var epListView = self.props.endpoints.map(function(ep){
				return (
					React.createElement("tr", {key: ep.id, className: "info"}, 
						React.createElement("td", null, ep.homingHost), 
	                    React.createElement("td", null, ep.contName), 
	                    React.createElement("td", null, ep.netID), 
						React.createElement("td", null, ep.ipAddress)
					)
				);
			});

			// Render the pane
			return (
	        React.createElement("div", {style: {margin: '5%',}}, 
				React.createElement(Table, {hover: true}, 
					React.createElement("thead", null, 
						React.createElement("tr", null, 
							React.createElement("th", null, "Host"), 
	                        React.createElement("th", null, "Container"), 
							React.createElement("th", null, "Network"), 
							React.createElement("th", null, "IP address")
						)
					), 
					React.createElement("tbody", null, 
	            		epListView
					)
				)
	        )
	    );
		}
	});

	module.exports = HomePane


/***/ },
/* 3 */
/***/ function(module, exports, __webpack_require__) {

	/** @jsx React.DOM */// network.js
	// Display Network information

	var NetworkPane = React.createClass({displayName: "NetworkPane",
	  	render: function() {
			var self = this

			if (self.props.networks === undefined) {
				return React.createElement("div", null, " ")
			}

			// Walk thru all the altas and see which ones are on this node
			var netListView = self.props.networks.map(function(network){
				if (network.isPublic) {
					netType = "public"
				} else {
					netType = "private"
				}
				return (
					React.createElement("tr", {key: network.key, className: "info"}, 
						React.createElement("td", null, network.tenantName), 
						React.createElement("td", null, network.networkName), 
						React.createElement("td", null, netType), 
						React.createElement("td", null, network.encap), 
						React.createElement("td", null, network.subnet)
					)
				);
			});

			// Render the pane
			return (
	        React.createElement("div", {style: {margin: '5%',}}, 
				React.createElement(Table, {hover: true}, 
					React.createElement("thead", null, 
						React.createElement("tr", null, 
							React.createElement("th", null, "Tenant"), 
							React.createElement("th", null, "Network"), 
							React.createElement("th", null, "Type"), 
							React.createElement("th", null, "Encapsulation"), 
							React.createElement("th", null, "Subnet")
						)
					), 
					React.createElement("tbody", null, 
	            		netListView
					)
				)
	        )
	    );
		}
	});

	module.exports = NetworkPane


/***/ },
/* 4 */
/***/ function(module, exports, __webpack_require__) {

	/** @jsx React.DOM */// groups.js
	// Display Endpoint group information

	var GroupsPane = React.createClass({displayName: "GroupsPane",
	  	render: function() {
			var self = this

			if (self.props.endpointGroups === undefined) {
				return React.createElement("div", null, " ")
			}

			// Walk thru all the altas and see which ones are on this node
			var epgListView = self.props.endpointGroups.map(function(epg){
				return (
					React.createElement("tr", {key: epg.key, className: "info"}, 
						React.createElement("td", null, epg.tenantName), 
						React.createElement("td", null, epg.networkName), 
						React.createElement("td", null, epg.groupName), 
						React.createElement("td", null, epg.policies)
					)
				);
			});

			// Render the pane
			return (
	        React.createElement("div", {style: {margin: '5%',}}, 
				React.createElement(Table, {hover: true}, 
					React.createElement("thead", null, 
						React.createElement("tr", null, 
							React.createElement("th", null, "Tenant"), 
							React.createElement("th", null, "Network"), 
							React.createElement("th", null, "Endpoint Group"), 
							React.createElement("th", null, "Policies")
						)
					), 
					React.createElement("tbody", null, 
	            		epgListView
					)
				)
	        )
	    );
		}
	});

	module.exports = GroupsPane


/***/ },
/* 5 */
/***/ function(module, exports, __webpack_require__) {

	/** @jsx React.DOM */// policy.js
	// Display Policy information

	var PolicyPane = React.createClass({displayName: "PolicyPane",
	  	render: function() {
			var self = this

			if (self.props.policies === undefined) {
				return React.createElement("div", null, " ")
			}

			// Walk thru all the altas and see which ones are on this node
			var policyListView = self.props.policies.map(function(policy){
				return (
					React.createElement("tr", {key: policy.key, className: "info"}, 
						React.createElement("td", null, policy.tenantName), 
						React.createElement("td", null, policy.policyName)
					)
				);
			});

			// Render the pane
			return (
	        React.createElement("div", {style: {margin: '5%',}}, 
				React.createElement(Table, {hover: true}, 
					React.createElement("thead", null, 
						React.createElement("tr", null, 
							React.createElement("th", null, "Tenant"), 
							React.createElement("th", null, "Policy")
						)
					), 
					React.createElement("tbody", null, 
	            		policyListView
					)
				)
	        )
	    );
		}
	});

	module.exports = PolicyPane


/***/ },
/* 6 */
/***/ function(module, exports, __webpack_require__) {

	/** @jsx React.DOM */// volumes.js
	// Display Volumes information

	var VolumesPane = React.createClass({displayName: "VolumesPane",
	  	render: function() {
			var self = this

			if (self.props.volumes === undefined) {
				return React.createElement("div", null, " ")
			}

			// Walk thru all the volumes
			var volListView = self.props.volumes.map(function(vol){
				return (
					React.createElement("tr", {key: vol.key, className: "info"}, 
						React.createElement("td", null, vol.tenantName), 
						React.createElement("td", null, vol.volumeName), 
						React.createElement("td", null, vol.poolName), 
						React.createElement("td", null, vol.size)
					)
				);
			});

			// Render the pane
			return (
	        React.createElement("div", {style: {margin: '5%',}}, 
				React.createElement(Table, {hover: true}, 
					React.createElement("thead", null, 
						React.createElement("tr", null, 
							React.createElement("th", null, "Tenant"), 
							React.createElement("th", null, "Volume"), 
							React.createElement("th", null, "Pool"), 
							React.createElement("th", null, "Size")
						)
					), 
					React.createElement("tbody", null, 
	            		volListView
					)
				)
	        )
	    );
		}
	});

	module.exports = VolumesPane


/***/ }
/******/ ]);
