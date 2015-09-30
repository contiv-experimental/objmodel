// policy.js
// Display Policy information

var PolicyPane = React.createClass({
  	render: function() {
		var self = this

		if (self.props.policies === undefined) {
			return <div> </div>
		}

		// Walk thru all the altas and see which ones are on this node
		var policyListView = self.props.policies.map(function(policy){
			return (
				<tr key={policy.key} className="info">
					<td>{policy.tenantName}</td>
					<td>{policy.policyName}</td>
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
						<th>Policy</th>
					</tr>
				</thead>
				<tbody>
            		{policyListView}
				</tbody>
			</Table>
        </div>
    );
	}
});

module.exports = PolicyPane
