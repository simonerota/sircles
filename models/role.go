package models

type RoleType string

// Don't change the names since these values are usually saved in the
// database
const (
	RoleTypeUndefined   		RoleType = "undefined"
	RoleTypeNormal      		RoleType = "normal"
	RoleTypeCircle      		RoleType = "circle"
	RoleTypeLeadLink    		RoleType = "leadlink"
	RoleTypeRepLink     		RoleType = "replink"
	RoleTypeFacilitator 		RoleType = "facilitator"
	RoleTypeEngager     		RoleType = "engager"
	RoleTypeChampion    		RoleType = "champion"
	RoleTypeScout       		RoleType = "scout"
	RoleTypeMagister    		RoleType = "magister"
	RoleTypeMangler     		RoleType = "mangler"
	RoleTypeSecretary   		RoleType = "secretary"
	RoleTypeSecurityEnabler		RoleType = "securityenabler"
	RoleTypeReporter	   		RoleType = "reporter"
)

func (r RoleType) IsCoreRoleType() bool {
	return r == RoleTypeLeadLink ||
		r == RoleTypeRepLink ||
		r == RoleTypeFacilitator ||
		r == RoleTypeEngager ||
		r == RoleTypeChampion ||
		r == RoleTypeScout ||
		r == RoleTypeMagister ||
		r == RoleTypeMangler ||
		r == RoleTypeSecretary ||
		r == RoleTypeSecurityEnabler ||
		r == RoleTypeReporter
}

func (r RoleType) String() string {
	return string(r)
}

func RoleTypeFromString(r string) RoleType {
	switch r {
	case "normal":
		return RoleTypeNormal
	case "circle":
		return RoleTypeCircle
	case "leadlink":
		return RoleTypeLeadLink
	case "replink":
		return RoleTypeRepLink
	case "facilitator":
		return RoleTypeFacilitator
	case "engager":
		return RoleTypeEngager
	case "champion":
		return RoleTypeChampion
	case "scout":
		return RoleTypeScout
	case "magister":
		return RoleTypeMagister
	case "mangler":
		return RoleTypeMangler
	case "secretary":
		return RoleTypeSecretary
	case "securityenabler":
		return RoleTypeSecurityEnabler
	case "reporter":
		return RoleTypeReporter
	default:
		return RoleTypeUndefined
	}
}

type Role struct {
	Vertex
	RoleType RoleType
	Depth    int32
	Name     string
	Purpose  string
}

type Roles []*Role

func (r Roles) Len() int           { return len(r) }
func (r Roles) Less(i, j int) bool { return r[i].Name < r[j].Name }
func (r Roles) Swap(i, j int)      { r[i], r[j] = r[j], r[i] }

type RoleAdditionalContent struct {
	Vertex
	Content string
}

type MemberCirclePermissions struct {
	AssignChildCircleLeadLink   bool
	AssignChildRoleMembers      bool
	ManageChildRoles            bool
	AssignCircleDirectMembers   bool
	AssignCircleCoreRoles       bool
	ManageRoleAdditionalContent bool
	// special cases for root circle
	AssignRootCircleLeadLink bool
	ManageRootCircle         bool
}

type CoreRoleDefinition struct {
	Role             *Role
	Domains          []*Domain
	Accountabilities []*Accountability
}

func GetCoreRoles() []*CoreRoleDefinition {
	return []*CoreRoleDefinition{
		{
			Role: &Role{
				Name:     "Sircle Leader",
				RoleType: RoleTypeLeadLink,
				Purpose:  "Sircle Leader holds the Purpose of the overall Sircle",
			},
			Domains: []*Domain{
				{Description: "Allocate Roles within the Sircle"},
			},
			Accountabilities: []*Accountability{
				{Description: "Structure Sircle’s Governance to enact its Purpose and Accountabilities"},
				{Description: "Sircle's Act Owner"},
				{Description: "Create new roles if needed"},
				{Description: "Allocate the Sircle's resources across its various Projects and/or Roles"},
				{Description: "Assign people to Sircle's Roles and monitor the fit"},
				{Description: "Offer feedback to enhance fit and re-assign Roles to other people when it could be useful for enhancing fit"},
				{Description: "Establish priorities and strategies for the Sircle"},
				{Description: "Define a general strategy for the Sircle, or multiple strategies, which are heuristics that guide the Sircle's roles in self-identifying priorities on an ongoing basis"},
				{Description: "Define kpi for the Sircle’s roles"},
				{Description: "Remove constraints and impediments within the Sircle"},
				{Description: "Manage conflicts"},
				{Description: "Identify the need of new sub-sircles and start the sub-sircle creation process with “Sircle of Staffing”"},
				{Description: "Team coaching"},
				{Description: "Write and keep up to date Value Propositions (leanvas.sorint.it)"},
				{Description: "Support “People Development Sircle” caring members through leveraging Sorint values"},
				{Description: "Owner of the staff capacity process. He reports any issues to “Sircle of Staffing”"},
				{Description: "Leader of the onboarding process of new core members"},
				{Description: "Owner of the performances evaluation process"},
				{Description: "Report and collaborate with “People Development Sircle” measuring core members performance behavior, motivation and skills improvement"},
			},
		},
		{
			Role: &Role{
				Name:     "Rep Link",
				RoleType: RoleTypeRepLink,
				Purpose:  "Within the Super-Circle, the Rep Link holds the Purpose of the SubCircle; within the Sub-Circle, the Rep Link’s Purpose is: Tensions relevant to process in the Super-Circle channeled out and resolved",
			},
			Accountabilities: []*Accountability{
				{Description: "Removing constraints within the broader Organization that limit the Sub-Circle"},
				{Description: "Seeking to understand Tensions conveyed by Sub-Circle Circle Members, and discerning those appropriate to process in the Super-Circle"},
				{Description: "Providing visibility to the Super-Circle into the health of the Sub-Circle, including reporting on any metrics or checklist items assigned to the whole Sub-Circle"},
			},
		},
		{
			Role: &Role{
				Name:     "Facilitator",
				RoleType: RoleTypeFacilitator,
				Purpose:  "Circle governance and operational practices aligned with the Constitution",
			},
			Accountabilities: []*Accountability{
				{Description: "Facilitating the Circle’s constitutionally-required meetings"},
				{Description: "Auditing the meetings and records of Sub-Circles as needed, and declaring a Process Breakdown upon discovering a pattern of behavior that conflicts with the rules of the Constitution"},
			},
		},
		{
			Role: &Role{
				Name:     "Engager",
				RoleType: RoleTypeEngager,
				Purpose:  "Extreme customer satisfaction proposing real added value solutions with visible Sorint signature",
			},
			Domains: []*Domain{
				{Description: "Technical Sircles"},
			},
			Accountabilities: []*Accountability{
				{Description: "Collect needs, pains and requirements directly from the customer to better understand the opportunities"},
				{Description: "Identify the right solution that matches business outcom"},
				{Description: "Provide the Statement of work for the activity which allows Sorint to build a trade offer defining solution studies (proposal, solution design, etc ...)"},
				{Description: "Define the number of resources needed to deliver activities and the required skills for each of them"},
				{Description: "At the end of all activities, provide a document useful for use cases and value proposition improvement"},
				{Description: "Present Sorint approach, methodologies, best practices…"},
				{Description: "Solution/Service/Product Presentation"},
				{Description: "Build and run Poc"},
				{Description: "Show Different Technologies through Use Cases"},
				{Description: "Help the Sircle Leader to update the Sircle Purpose, Domain and Value proposition after every Strategy Meetings"},
			},
		},
		{
			Role: &Role{
				Name:     "Champion",
				RoleType: RoleTypeChampion,
				Purpose:  "Extreme customer satisfaction and establishing valued relationship with visible Sorint signature, the Sorint ambassador",
			},
			Domains: []*Domain{
				{Description: "Customer focus Sircle"},
			},
			Accountabilities: []*Accountability{
				{Description: "Identify and map customer organization and various managers"},
				{Description: "Study the technical sircles's value propositions"},
				{Description: "Arrange regular meetings with the customer key people in order to share the vp"},
				{Description: "Align with the Sircle leader for opportunities"},
				{Description: "Open opportunities for any intercepted customer need"},
				{Description: "Report competitors presence on the customer site (name and where)"},
				{Description: "Report negative and positive customer feedback on our services"},
				{Description: "Report negative and positive customer feedback on competitors works"},
				{Description: "Report any changes on the site that could lead to opportunities, instability or problems"},
				{Description: "Report customer intentions to purchase new products"},
				{Description: "Be aware of all customer-related opportunities"},
				{Description: "Organize meetings/phone calls with Sorint Business Developers for alignment"},
				{Description: "Help the Sircle Leader to update the Sircle Purpose, Domain and Value proposition after every Strategy Meeting"},
			},
		},		
		{
			Role: &Role{
				Name:     "Talent Handler",
				RoleType: RoleTypeScout,
				Purpose:  "Help Sorintians to have only talented collegues",
			},
			Domains: []*Domain{
				{Description: "Sircle"},
			},
			Accountabilities: []*Accountability{
				{Description: "Responsible for  keeping core members CV updated (My BIO)"},
				{Description: "Understand the right characteristics of the ideal candidate for customers requests"},
				{Description: "Help “Sircle of Staffing” to review the candidate’s CV"},
				{Description: "Help “Sircle of Staffing” or “People Development  Sircle” with the candidates interviews"},
				{Description: "Help candidates to improve CV and Linkedin in order to give the right information to the customer"},
				{Description: "Give the candidate all the needed information in order to cope with the customer interview(dressing code, key points, suggestions …)"},
				{Description: "Introduce the candidate to the customer (sending CV, attending the meeting)"},
			},
		},
		{
			Role: &Role{
				Name:     "Academier",
				RoleType: RoleTypeMagister,
				Purpose:  "Neverending Sorintians Growth",
			},
			Domains: []*Domain{
				{Description: "Sircle Core Members"},
			},
			Accountabilities: []*Accountability{
				{Description: "Following the Sircle Leader growth strategies identify skill gaps"},
				{Description: "Write a Training Path"},
				{Description: "Identify course in academy hub"},
				{Description: "Suggest Academia about courses in the market"},
				{Description: "Keep aligned Academia on the released trainings pat"},
				{Description: "Monitor the educational path results"},
				{Description: "Responsible for the apprentices' training plan"},
			},
		},
		{
			Role: &Role{
				Name:     "Planner",
				RoleType: RoleTypeMangler,
				Purpose:  "The right project in the right time with maximum satisfaction",
			},
			Domains: []*Domain{
				{Description: "Services Activities under Sircle's responsibility"},
			},
			Accountabilities: []*Accountability{
				{Description: "Schedule properly assigning the right member to the right “ACO Profile” in agreement with the Sircle Leader’s instructions"},
				{Description: "Keep the schedule update"},
				{Description: "Schedule properly in order to successfully complete project/services (on time, with the right effort)"},
				{Description: "Manage resources requests conflicts"},
				{Description: "Inform PMO about people Availability"},
				{Description: "Keep all stakeholders updated about projects or services progress"},
				{Description: "Pay attention and act in advance planning critical periods (holidays …)"},
				{Description: "Holidays disposal"},
				{Description: "Solve NC"},
				{Description: "Close ACOs as soon as the project is finished"},
			},
		},
		{
			Role: &Role{
				Name:     "Secretary",
				RoleType: RoleTypeSecretary,
				Purpose:  "Steward and stabilize the Circle’s formal records and record-keeping process",
			},
			Domains: []*Domain{
				{Description: "All constitutionally-required records of the Circle"},
			},
			Accountabilities: []*Accountability{
				{Description: "Scheduling the Circle’s required meetings, and notifying all Core Circle Members of scheduled times and locations"},
				{Description: "Capturing and publishing the outputs of the Circle’s required meetings, and maintaining a compiled view of the Circle’s current Governance, checklist items, and metrics"},
				{Description: "Interpreting Governance and the Constitution upon request"},
			},
		},
		{
			Role: &Role{
				Name:     "Security Enabler",
				RoleType: RoleTypeSecurityEnabler,
				Purpose:  "Ensure that personal data, Sorint data and customer data are protected",
			},
			Domains: []*Domain{
				{Description: "Sircle"},
			},
			Accountabilities: []*Accountability{
				{Description: "Know Sorint policies and procedures on security(https://clann.sorint.it/policies-procedures2)"},
				{Description: "Ensure the information security (data, communications, documents...) based on the classification"},
				{Description: "Contribute to reporting improvements to Sorint best practices"},
			},
		},
		{
			Role: &Role{
				Name:     "Reporter",
				RoleType: RoleTypeReporter,
				Purpose:  "Make all Sorint participate in what happens in the Sircle and make the Sircle participate in what happens in Sorint",
			},
			Domains: []*Domain{
				{Description: "Sircle"},
			},
			Accountabilities: []*Accountability{
				{Description: "Let Sorint's activities be internally perceived: events, prj, news, various updates…"},
				{Description: "Direct core members to inquire about what happens in the other Sircles"},
				{Description: "Promote in the Sorint Sircle posts published on social networks"},
				{Description: "Advertise business events"},
				{Description: "Publish a post a monthtelling what happened in the Sircle (success stories, news, projects evolution, courses…)"},
			},
		},
	}
}
