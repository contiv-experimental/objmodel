// network.js
// Display Network information

var NetworkPane = React.createClass({
  	render: function() {
		var self = this

		if (self.props.networks === undefined) {
			return <div> </div>
		}

		// Walk thru all the altas and see which ones are on this node
		var netListView = self.props.networks.map(function(network){
			if (network.isPublic) {
				netType = "public"
			} else {
				netType = "private"
			}
			return (
				<tr key={network.key} className="info">
					<td>{network.tenantName}</td>
					<td>{network.networkName}</td>
					<td>{netType}</td>
					<td>{network.encap}</td>
					<td>{network.subnet}</td>
				</tr>
			);
		});

		// Render the pane
		return (
        <div style={{margin: '5%',}}>
			<Table hover>
				<thead>
					<tr>
						<th>Tenant</th>
						<th>Network</th>
						<th>Type</th>
						<th>Encapsulation</th>
						<th>Subnet</th>
					</tr>
				</thead>
				<tbody>
            		{netListView}
				</tbody>
			</Table>
        </div>
    );
	}
});

module.exports = NetworkPane
