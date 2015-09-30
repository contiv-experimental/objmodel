// groups.js
// Display Endpoint group information

var GroupsPane = React.createClass({
  	render: function() {
		var self = this

		if (self.props.endpointGroups === undefined) {
			return <div> </div>
		}

		// Walk thru all the altas and see which ones are on this node
		var epgListView = self.props.endpointGroups.map(function(epg){
			return (
				<tr key={epg.key} className="info">
					<td>{epg.tenantName}</td>
					<td>{epg.networkName}</td>
					<td>{epg.groupName}</td>
					<td>{epg.policies}</td>
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
						<th>Endpoint Group</th>
						<th>Policies</th>
					</tr>
				</thead>
				<tbody>
            		{epgListView}
				</tbody>
			</Table>
        </div>
    );
	}
});

module.exports = GroupsPane
