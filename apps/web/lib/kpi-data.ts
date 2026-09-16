import { KPIDefinition, LiveFeedTick, MacroKPISummary, Project, ProjectKPIScorecard, ProjectMetricAssessment, PillarScore } from "./types";

export const CANONICAL_KPIS: KPIDefinition[] = [
  // Pillar 1: ESG & Decarbonization
  {
    code: "ESG.GHG.INTENSITY.SCOPE1_2",
    name: "Scope 1 & 2 Emissions Intensity",
    category: "ESG_DECARBONIZATION",
    description: "Direct and indirect operational greenhouse gas emissions per $1M CAD capex or per tonne output.",
    unit: "tCO2e/$1M CAD",
    preferred_direction: "LOWER_BETTER",
    target_benchmark: 45.0,
    benchmark_unit: "tCO2e/$1M CAD",
    benchmark_label: "Net-Zero 2050 Sectoral Pathway Threshold",
    statutory_basis: "Canadian Net-Zero Emissions Accountability Act",
    publisher: "Environment and Climate Change Canada",
    update_frequency: "ANNUAL",
    is_live_feed: false,
  },
  {
    code: "ESG.GHG.ABATEMENT.ANNUAL_MT",
    name: "Annual Lifecycle GHG Abatement",
    category: "ESG_DECARBONIZATION",
    description: "Net greenhouse gas emissions avoided or sequestered annually relative to baseline fossil counter-factual.",
    unit: "Mt CO2e/yr",
    preferred_direction: "HIGHER_BETTER",
    target_benchmark: 1.5,
    benchmark_unit: "Mt CO2e/yr",
    benchmark_label: "Canada 2030 Emissions Reduction Plan Contributor",
    statutory_basis: "Federal 2030 Emissions Reduction Plan (ERP)",
    publisher: "Clean Energy Canada / ECCC National Inventory",
    update_frequency: "ANNUAL",
    is_live_feed: false,
  },
  {
    code: "ESG.CARBON.TAX_SENSITIVITY.CAD",
    name: "OBPS Carbon Tax Exposure Sensitivity",
    category: "ESG_DECARBONIZATION",
    description: "Marginal annual project operating cost impact per $10/tonne escalation in Federal OBPS benchmark carbon price.",
    unit: "$CAD/tonne CO2e",
    preferred_direction: "LOWER_BETTER",
    target_benchmark: 170.0,
    benchmark_unit: "$CAD/tonne",
    benchmark_label: "2030 Statutory Carbon Price Cap ($170/t)",
    statutory_basis: "Greenhouse Gas Pollution Pricing Act (GGPPA)",
    publisher: "Department of Finance Canada / CRA",
    update_frequency: "ANNUAL",
    is_live_feed: false,
  },
  {
    code: "ESG.WATER.RECYCLING.RATIO",
    name: "Industrial Water Circularity & Recycling",
    category: "ESG_DECARBONIZATION",
    description: "Percentage of operational process water recycled or closed-loop circulated on site.",
    unit: "percent",
    preferred_direction: "HIGHER_BETTER",
    target_benchmark: 85.0,
    benchmark_unit: "percent",
    benchmark_label: "Best-in-Class Zero Liquid Discharge Standard",
    statutory_basis: "Fisheries Act / Mining Effluent Regulations",
    publisher: "ECCC Water Quality Monitoring",
    update_frequency: "QUARTERLY",
    is_live_feed: false,
  },
  {
    code: "ESG.POWER.CLEAN_PURITY.PCT",
    name: "Clean Zero-Carbon Power Utilization",
    category: "ESG_DECARBONIZATION",
    description: "Percentage of connected electricity consumption supplied from zero-emitting generation.",
    unit: "percent",
    preferred_direction: "HIGHER_BETTER",
    target_benchmark: 95.0,
    benchmark_unit: "percent",
    benchmark_label: "Clean Electricity Regulations 2035 Compliance",
    statutory_basis: "Federal Clean Electricity Regulations (CER)",
    publisher: "Provincial System Operators (IESO, AESO, HQ, BC Hydro)",
    update_frequency: "MONTHLY",
    is_live_feed: false,
  },

  // Pillar 2: Indigenous Economic Sovereignty
  {
    code: "INDIG.EQUITY.OWNERSHIP.PCT",
    name: "First Nations Equity Co-Ownership Stake",
    category: "INDIGENOUS_EQUITY",
    description: "Percentage of project commercial equity owned by host First Nations, Inuit, or Métis communities.",
    unit: "percent",
    preferred_direction: "HIGHER_BETTER",
    target_benchmark: 25.0,
    benchmark_unit: "percent",
    benchmark_label: "CIB / National Indigenous Economic Strategy Target",
    statutory_basis: "UNDRIP Act & TRC Call to Action 92",
    publisher: "Canada Infrastructure Bank / FNMPC",
    update_frequency: "EVENT_DRIVEN",
    is_live_feed: false,
  },
  {
    code: "INDIG.PROCUREMENT.SHARE.PCT",
    name: "Indigenous Enterprise Procurement Share",
    category: "INDIGENOUS_EQUITY",
    description: "Percentage of Tier-1 and Tier-2 procurement contracts awarded to verified Indigenous-owned businesses.",
    unit: "percent",
    preferred_direction: "HIGHER_BETTER",
    target_benchmark: 10.0,
    benchmark_unit: "percent",
    benchmark_label: "Exceeds Federal 5% Mandatory Baseline",
    statutory_basis: "Directive on Government Contracts Awarded to Indigenous Businesses",
    publisher: "Indigenous Services Canada / CCAB",
    update_frequency: "QUARTERLY",
    is_live_feed: false,
  },
  {
    code: "INDIG.WORKFORCE.PARTICIPATION.PCT",
    name: "Local Indigenous Workforce Participation",
    category: "INDIGENOUS_EQUITY",
    description: "Proportion of peak construction and steady-state operating workforce from host Indigenous communities.",
    unit: "percent",
    preferred_direction: "HIGHER_BETTER",
    target_benchmark: 20.0,
    benchmark_unit: "percent",
    benchmark_label: "Northern Impact Benefit Agreement (IBA) Standard",
    statutory_basis: "Impact Benefit Agreements / Section 35",
    publisher: "StatCan Indigenous Peoples Survey & Project Disclosures",
    update_frequency: "ANNUAL",
    is_live_feed: false,
  },
  {
    code: "INDIG.ILGP.GUARANTEE.CAD",
    name: "Indigenous Loan Guarantee Program Allocation",
    category: "INDIGENOUS_EQUITY",
    description: "Federal debt guarantee allocated under the $5B Canada Indigenous Loan Guarantee Program.",
    unit: "CAD",
    preferred_direction: "HIGHER_BETTER",
    target_benchmark: 100_000_000,
    benchmark_unit: "$CAD",
    benchmark_label: "$5B Federal Pool Allocation Capacity",
    statutory_basis: "Budget 2024 / NRCan ILGP Program",
    publisher: "Natural Resources Canada / CDEV",
    update_frequency: "EVENT_DRIVEN",
    is_live_feed: false,
  },
  {
    code: "INDIG.DIVIDEND.COMMUNITY.CAD_YR",
    name: "30-Year Annualized Community Dividend Yield",
    category: "INDIGENOUS_EQUITY",
    description: "Net untied sovereign annual cash flow distributed to community trust accounts for intergenerational prosperity.",
    unit: "CAD/yr",
    preferred_direction: "HIGHER_BETTER",
    target_benchmark: 12_000_000,
    benchmark_unit: "$CAD/yr",
    benchmark_label: "Sovereign Intergenerational Community Fund",
    statutory_basis: "First Nations Fiscal Management Act",
    publisher: "First Nations Financial Management Board",
    update_frequency: "ANNUAL",
    is_live_feed: false,
  },

  // Pillar 3: Capital Spend Velocity & Hazards
  {
    code: "CAPEX.SPEND.VELOCITY.CAD_MO",
    name: "Capital Deployment Spend Velocity",
    category: "CAPITAL_VELOCITY",
    description: "Average monthly capital deployment and procurement cash-out run-rate during active execution.",
    unit: "CAD_M/month",
    preferred_direction: "HIGHER_BETTER",
    target_benchmark: 35.0,
    benchmark_unit: "$M CAD/mo",
    benchmark_label: "Efficient Megaproject Spend Rate Baseline",
    statutory_basis: "Major Projects Management Office Reporting",
    publisher: "Corporate SEDAR+ MD&A",
    update_frequency: "MONTHLY",
    is_live_feed: false,
  },
  {
    code: "SCHEDULE.SLIP.MONTHS",
    name: "Cumulative Critical-Path Schedule Slip",
    category: "CAPITAL_VELOCITY",
    description: "Months variance between current forecasted commercial operation date (COD) and baseline FID target.",
    unit: "months",
    preferred_direction: "LOWER_BETTER",
    target_benchmark: 0.0,
    benchmark_unit: "months",
    benchmark_label: "Zero-Slippage Critical Path Delivery",
    statutory_basis: "Treasury Board Directive on Management of Major Projects",
    publisher: "Regulatory Filings / Independent Engineering Audits",
    update_frequency: "QUARTERLY",
    is_live_feed: false,
  },
  {
    code: "RISK.COST_OVERRUN.HAZARD.PCT",
    name: "Flyvbjerg Cost Overrun Hazard Probability",
    category: "CAPITAL_VELOCITY",
    description: "Empirical Bayesian lognormal probability of exceeding approved capital budget by >20%.",
    unit: "percent",
    preferred_direction: "LOWER_BETTER",
    target_benchmark: 25.0,
    benchmark_unit: "percent",
    benchmark_label: "Low-Tail-Risk Benchmark (<25%)",
    statutory_basis: "Infrastructure Canada Capital Risk Guidelines",
    publisher: "Oxford Megaproject Database / CEO-G Risk Engine",
    update_frequency: "CONTINUAL",
    is_live_feed: false,
  },
  {
    code: "FINANCE.PUBLIC_CROWDING.MULTIPLIER",
    name: "Private Capital Crowding-In Multiplier",
    category: "CAPITAL_VELOCITY",
    description: "Ratio of private institutional capital mobilized per $1 CAD of federal/provincial concessionary debt or grant.",
    unit: "ratio_x",
    preferred_direction: "HIGHER_BETTER",
    target_benchmark: 4.0,
    benchmark_unit: "x",
    benchmark_label: "CIB Legislative Target (3x-5x Private Crowding-In)",
    statutory_basis: "Canada Infrastructure Bank Act",
    publisher: "Canada Infrastructure Bank Reports",
    update_frequency: "EVENT_DRIVEN",
    is_live_feed: false,
  },

  // Pillar 4: Sovereign Supply Chain & Domestic Content
  {
    code: "SUPPLY.DOMESTIC_CONTENT.PCT",
    name: "Canadian Domestic Content Share",
    category: "SUPPLY_CHAIN_CONTENT",
    description: "Proportion of direct capital expenditures procured from suppliers registered with a Canadian Business Number (BN).",
    unit: "percent",
    preferred_direction: "HIGHER_BETTER",
    target_benchmark: 70.0,
    benchmark_unit: "percent",
    benchmark_label: "Bill C-59 Domestic Industrial Benefit Target",
    statutory_basis: "Fall Economic Statement Implementation Act (Bill C-59)",
    publisher: "ISED Canadian Company Capabilities & CRA BN Registry",
    update_frequency: "SEMI_ANNUAL",
    is_live_feed: false,
  },
  {
    code: "SUPPLY.BOTTLENECK.LEAD_TIME.SCORE",
    name: "Critical Equipment Bottleneck Index",
    category: "SUPPLY_CHAIN_CONTENT",
    description: "Composite risk rating (1-100) assessing global procurement queues for transformers, calandria, turbines, and metallurgy.",
    unit: "score_1_to_100",
    preferred_direction: "LOWER_BETTER",
    target_benchmark: 30.0,
    benchmark_unit: "points",
    benchmark_label: "Low Supply-Chain Disruption Risk (<30)",
    statutory_basis: "National Supply Chain Strategy / Transport Canada",
    publisher: "Transport Canada Task Force & IMF PortWatch",
    update_frequency: "MONTHLY",
    is_live_feed: false,
  },
  {
    code: "TRADE.US_TARIFF.EXPOSURE.PCT",
    name: "U.S. Cross-Border Tariff Vulnerability",
    category: "SUPPLY_CHAIN_CONTENT",
    description: "Percentage of projected commercial offtake revenue subject to U.S. Section 232/301 or border adjustment trade friction.",
    unit: "percent",
    preferred_direction: "LOWER_BETTER",
    target_benchmark: 15.0,
    benchmark_unit: "percent",
    benchmark_label: "Diversified Global Trade Offtake Shield",
    statutory_basis: "Canada-United States-Mexico Agreement (CUSMA)",
    publisher: "Global Affairs Canada Trade Controls Bureau",
    update_frequency: "QUARTERLY",
    is_live_feed: false,
  },
  {
    code: "TRADE.IMPORT_REPLACEMENT.CAD_YR",
    name: "Annual Domestic Import Replacement Value",
    category: "SUPPLY_CHAIN_CONTENT",
    description: "Estimated annual dollar value of foreign imports directly displaced by project output.",
    unit: "CAD_M/yr",
    preferred_direction: "HIGHER_BETTER",
    target_benchmark: 250.0,
    benchmark_unit: "$M CAD/yr",
    benchmark_label: "Major Sovereign Industrial Contribution",
    statutory_basis: "Canadian Critical Minerals Strategy",
    publisher: "Statistics Canada Trade Data Online",
    update_frequency: "ANNUAL",
    is_live_feed: false,
  },

  // Pillar 5: Power, Grid Interconnection & Hosting
  {
    code: "GRID.PEAK_POWER.DEMAND.MW",
    name: "Connected Peak Electrical Load / Gen",
    category: "GRID_PHYSICS",
    description: "Total firm electrical load or dispatchable clean power capacity connected to the bulk electricity system.",
    unit: "MW",
    preferred_direction: "NEUTRAL",
    target_benchmark: 300.0,
    benchmark_unit: "MW",
    benchmark_label: "Transmission-Class Industrial Interconnection",
    statutory_basis: "National Energy Board / Electricity Acts",
    publisher: "Provincial ISOs / CER",
    update_frequency: "EVENT_DRIVEN",
    is_live_feed: false,
  },
  {
    code: "GRID.QUEUE.WAIT_TIME.MONTHS",
    name: "Grid Interconnection Queue Latency",
    category: "GRID_PHYSICS",
    description: "Elapsed months from initial System Impact Assessment (SIA) filing to signed Interconnection Agreement.",
    unit: "months",
    preferred_direction: "LOWER_BETTER",
    target_benchmark: 18.0,
    benchmark_unit: "months",
    benchmark_label: "Streamlined SIA Queue Latency",
    statutory_basis: "IESO / AESO Market Rules",
    publisher: "IESO / AESO Transmission Queue Registers",
    update_frequency: "MONTHLY",
    is_live_feed: false,
  },
  {
    code: "GRID.HOSTING.HEADROOM.MW",
    name: "Substation Electrical Hosting Headroom",
    category: "GRID_PHYSICS",
    description: "Available unallocated transformation capacity at nearest transmission node prior to requiring system upgrades.",
    unit: "MW",
    preferred_direction: "HIGHER_BETTER",
    target_benchmark: 150.0,
    benchmark_unit: "MW",
    benchmark_label: "Sufficient Local Transformation Margin",
    statutory_basis: "Transmission System Code",
    publisher: "Hydro One / AltaLink / BC Hydro Planning",
    update_frequency: "QUARTERLY",
    is_live_feed: false,
  },
  {
    code: "GRID.REINFORCEMENT.COST.CAD",
    name: "Upstream Network Reinforcement Cost",
    category: "GRID_PHYSICS",
    description: "Required external capital allocation to reinforce regional 230kV/500kV bulk transmission lines.",
    unit: "CAD",
    preferred_direction: "LOWER_BETTER",
    target_benchmark: 25_000_000,
    benchmark_unit: "$CAD",
    benchmark_label: "Minimal Network Reinforcement Burden",
    statutory_basis: "Provincial Utility Rate Hearings",
    publisher: "Independent System Operator Assessments",
    update_frequency: "EVENT_DRIVEN",
    is_live_feed: false,
  },

  // Pillar 6: Regulatory, Permitting & Legal Latency
  {
    code: "REG.IAAC.REVIEW_DURATION.MONTHS",
    name: "Federal Impact Assessment Duration",
    category: "REGULATORY_SPEED",
    description: "Cumulative elapsed months in federal environmental impact review from description to Decision Statement.",
    unit: "months",
    preferred_direction: "LOWER_BETTER",
    target_benchmark: 24.0,
    benchmark_unit: "months",
    benchmark_label: "Bill C-69 Statutory 2-Year Target",
    statutory_basis: "Impact Assessment Act (S.C. 2019, c. 28)",
    publisher: "Impact Assessment Agency of Canada (IAAC)",
    update_frequency: "MONTHLY",
    is_live_feed: false,
  },
  {
    code: "REG.PERMIT.COMPLETION.PCT",
    name: "Statutory Permitting Completion Progress",
    category: "REGULATORY_SPEED",
    description: "Percentage of required municipal, provincial water/mining, and federal approvals formally executed.",
    unit: "percent",
    preferred_direction: "HIGHER_BETTER",
    target_benchmark: 90.0,
    benchmark_unit: "percent",
    benchmark_label: "Construction-Ready Statutory Threshold",
    statutory_basis: "Provincial Environmental Protection Acts",
    publisher: "Provincial EA Registries (BC EAO, ERO, AER)",
    update_frequency: "MONTHLY",
    is_live_feed: false,
  },
  {
    code: "REG.LITIGATION.RISK_INDEX",
    name: "Litigious Friction & Judicial Review Risk",
    category: "REGULATORY_SPEED",
    description: "Quantitative exposure score (1-100) assessing probability of Section 35 constitutional challenges.",
    unit: "score_1_to_100",
    preferred_direction: "LOWER_BETTER",
    target_benchmark: 20.0,
    benchmark_unit: "points",
    benchmark_label: "Low Judicial Injunction Risk (<20)",
    statutory_basis: "Constitution Act, 1982, Section 35",
    publisher: "Federal Court of Canada / Legal Registries",
    update_frequency: "CONTINUAL",
    is_live_feed: false,
  },
  {
    code: "REG.CONDITIONS.LEGALLY_BINDING.COUNT",
    name: "Binding Environmental Certificate Conditions",
    category: "REGULATORY_SPEED",
    description: "Total count of legally enforceable conditions attached to Environmental Assessment Approval Certificate.",
    unit: "count",
    preferred_direction: "NEUTRAL",
    target_benchmark: 65.0,
    benchmark_unit: "conditions",
    benchmark_label: "Manageable Compliance Envelope",
    statutory_basis: "IAAC Decision Statement Terms",
    publisher: "IAAC / Provincial EA Registries",
    update_frequency: "EVENT_DRIVEN",
    is_live_feed: false,
  },

  // Pillar 7: Labour Dynamics & Critical Skills Pipeline
  {
    code: "LABOR.RED_SEAL.TRADES_GAP.FTE",
    name: "Peak Construction Red Seal Trades Deficit",
    category: "LABOR_SKILLS",
    description: "Estimated regional shortfall of certified union journeypersons at peak construction.",
    unit: "FTE",
    preferred_direction: "LOWER_BETTER",
    target_benchmark: 0.0,
    benchmark_unit: "FTE",
    benchmark_label: "Fully Balanced Craft Labour Supply",
    statutory_basis: "Red Seal Program / ESDC",
    publisher: "BuildForce Canada / CBTU",
    update_frequency: "ANNUAL",
    is_live_feed: false,
  },
  {
    code: "LABOR.APPRENTICE.RATIO.PCT",
    name: "Apprentice-to-Journeyperson Ratio",
    category: "LABOR_SKILLS",
    description: "Proportion of hours worked on job sites by registered apprentices (Bill C-59 prevailing wage bonus compliance).",
    unit: "percent",
    preferred_direction: "HIGHER_BETTER",
    target_benchmark: 12.0,
    benchmark_unit: "percent",
    benchmark_label: "Bill C-59 Clean Economy ITC Bonus Threshold (>=10%)",
    statutory_basis: "Income Tax Act s. 127.44",
    publisher: "CRA / Provincial Apprenticeship Boards",
    update_frequency: "QUARTERLY",
    is_live_feed: false,
  },
  {
    code: "LABOR.HOUSING.ABSORPTION_DEFICIT.UNITS",
    name: "Local Housing Absorption Deficit",
    category: "LABOR_SKILLS",
    description: "Shortfall of rental housing units within 45 min commute radius to accommodate inbound construction workforce.",
    unit: "units",
    preferred_direction: "LOWER_BETTER",
    target_benchmark: 0.0,
    benchmark_unit: "units",
    benchmark_label: "Adequate Municipal Rental Headroom",
    statutory_basis: "National Housing Strategy Act",
    publisher: "Canada Mortgage and Housing Corporation (CMHC)",
    update_frequency: "SEMI_ANNUAL",
    is_live_feed: false,
  },
  {
    code: "LABOR.TOTAL.PERSON_YEARS.COUNT",
    name: "Total Lifecycle Person-Years Created",
    category: "LABOR_SKILLS",
    description: "Combined direct, indirect supply chain, and induced consumer spending employment created.",
    unit: "person_years",
    preferred_direction: "HIGHER_BETTER",
    target_benchmark: 8_500.0,
    benchmark_unit: "person_years",
    benchmark_label: "Major Employment Driver",
    statutory_basis: "StatCan Input-Output Multipliers",
    publisher: "Statistics Canada SUT Model",
    update_frequency: "EVENT_DRIVEN",
    is_live_feed: false,
  },

  // Pillar 8: Live Commodity & Macro Benchmark Feeds
  {
    code: "COMMODITY.WCS_WTI.DIFF.USD_BBL",
    name: "Western Canadian Select (WCS) Discount",
    category: "COMMODITY_MACRO",
    description: "Heavy crude pricing differential at Hardisty, Alberta relative to Cushing WTI.",
    unit: "USD/bbl",
    preferred_direction: "LOWER_BETTER",
    target_benchmark: 12.50,
    benchmark_unit: "USD/bbl",
    benchmark_label: "Narrow Discount Baseline (<$13.00)",
    statutory_basis: "Canada Energy Regulator Pricing",
    publisher: "ICE Crude / Argus Hardisty",
    update_frequency: "REALTIME",
    is_live_feed: true,
  },
  {
    code: "COMMODITY.AECO.GAS.CAD_GJ",
    name: "AECO C Natural Gas Spot Price",
    category: "COMMODITY_MACRO",
    description: "Western Canadian benchmark natural gas hub price at NIT Nova Inventory Transfer.",
    unit: "CAD/GJ",
    preferred_direction: "NEUTRAL",
    target_benchmark: 2.75,
    benchmark_unit: "CAD/GJ",
    benchmark_label: "Economic Feedstock Parity",
    statutory_basis: "AUC / CER Gas Pricing",
    publisher: "NGX Alberta NIT Hub",
    update_frequency: "DAILY",
    is_live_feed: true,
  },
  {
    code: "COMMODITY.LME.NICKEL.USD_T",
    name: "LME Grade-1 Nickel Cash Price",
    category: "COMMODITY_MACRO",
    description: "London Metal Exchange cash contract for primary nickel for battery cathode offtake.",
    unit: "USD/tonne",
    preferred_direction: "HIGHER_BETTER",
    target_benchmark: 18_500.0,
    benchmark_unit: "USD/t",
    benchmark_label: "Canadian SMR & EV Mine Breakeven",
    statutory_basis: "Critical Minerals Offtake Benchmarks",
    publisher: "London Metal Exchange (LME)",
    update_frequency: "REALTIME",
    is_live_feed: true,
  },
  {
    code: "COMMODITY.UX.U3O8.USD_LB",
    name: "Ux U3O8 Spot Uranium Price",
    category: "COMMODITY_MACRO",
    description: "Industry benchmark spot price per pound of triuranium octoxide for CANDU/SMR reactor fuel.",
    unit: "USD/lb",
    preferred_direction: "NEUTRAL",
    target_benchmark: 85.0,
    benchmark_unit: "USD/lb",
    benchmark_label: "Stable Long-Term Fuel Cycle Window",
    statutory_basis: "Nuclear Safety and Control Act",
    publisher: "UxC Weekly Spot Indicator",
    update_frequency: "WEEKLY",
    is_live_feed: true,
  },
  {
    code: "COMMODITY.LITHIUM.CARBONATE.USD_T",
    name: "Battery-Grade Lithium Carbonate",
    category: "COMMODITY_MACRO",
    description: "Spot benchmark for North American delivered battery precursor lithium carbonate.",
    unit: "USD/tonne",
    preferred_direction: "HIGHER_BETTER",
    target_benchmark: 14_000.0,
    benchmark_unit: "USD/t",
    benchmark_label: "Commercial Extraction Incentive Threshold",
    statutory_basis: "Federal Critical Minerals Infrastructure Fund",
    publisher: "Benchmark Minerals North America",
    update_frequency: "WEEKLY",
    is_live_feed: true,
  },
  {
    code: "MACRO.CAD_USD.SPOT_FX",
    name: "Bank of Canada CAD/USD Spot Rate",
    category: "COMMODITY_MACRO",
    description: "Official Bank of Canada daily noon fixing for Canadian Dollar against U.S. Dollar.",
    unit: "CAD_per_USD",
    preferred_direction: "NEUTRAL",
    target_benchmark: 0.735,
    benchmark_unit: "USD_per_CAD",
    benchmark_label: "Monetary Neutral Range",
    statutory_basis: "Bank of Canada Act",
    publisher: "Bank of Canada Valet Noon Fix",
    update_frequency: "BUSINESS_DAILY",
    is_live_feed: true,
  },
  {
    code: "MACRO.BOC.POLICY_RATE.PCT",
    name: "Bank of Canada Policy Interest Rate",
    category: "COMMODITY_MACRO",
    description: "Bank of Canada target for the overnight lending rate governing institutional debt facilities.",
    unit: "percent",
    preferred_direction: "NEUTRAL",
    target_benchmark: 2.75,
    benchmark_unit: "percent",
    benchmark_label: "Neutral Real Policy Rate",
    statutory_basis: "Bank of Canada Monetary Policy",
    publisher: "Bank of Canada Valet",
    update_frequency: "FIXED_DATES",
    is_live_feed: true,
  },
  {
    code: "MACRO.GOC.10Y_YIELD.PCT",
    name: "Government of Canada 10-Year Bond Yield",
    category: "COMMODITY_MACRO",
    description: "Yield to maturity on 10-year GoC sovereign bonds serving as risk-free discount rate.",
    unit: "percent",
    preferred_direction: "NEUTRAL",
    target_benchmark: 2.95,
    benchmark_unit: "percent",
    benchmark_label: "Project Finance Hurdle Base",
    statutory_basis: "Financial Administration Act",
    publisher: "Bank of Canada Valet Markets",
    update_frequency: "BUSINESS_DAILY",
    is_live_feed: true,
  },
];

export const INITIAL_LIVE_TICKS: LiveFeedTick[] = [
  { sequence_id: 101, metric_code: "COMMODITY.WCS_WTI.DIFF.USD_BBL", name: "WCS-WTI Discount", category: "COMMODITY_MACRO", value: 12.85, unit: "USD/bbl", change_absolute: -0.35, change_percent: -2.65, direction: "DOWN", source: "ICE / Argus Hardisty", timestamp: "2026-09-16T22:30:00Z", hash: "a8e94f1c" },
  { sequence_id: 102, metric_code: "COMMODITY.AECO.GAS.CAD_GJ", name: "AECO C Gas Spot", category: "COMMODITY_MACRO", value: 2.48, unit: "CAD/GJ", change_absolute: 0.08, change_percent: 3.33, direction: "UP", source: "NGX Alberta NIT Hub", timestamp: "2026-09-16T22:30:00Z", hash: "7c18fa3b" },
  { sequence_id: 103, metric_code: "COMMODITY.LME.NICKEL.USD_T", name: "LME Nickel Cash", category: "COMMODITY_MACRO", value: 18950.0, unit: "USD/t", change_absolute: 125.0, change_percent: 0.66, direction: "UP", source: "London Metal Exchange", timestamp: "2026-09-16T22:30:00Z", hash: "3e5b10da" },
  { sequence_id: 104, metric_code: "COMMODITY.UX.U3O8.USD_LB", name: "Ux U3O8 Uranium", category: "COMMODITY_MACRO", value: 86.50, unit: "USD/lb", change_absolute: 0.50, change_percent: 0.58, direction: "UP", source: "UxC Cameco Spot Index", timestamp: "2026-09-16T22:30:00Z", hash: "89cbf410" },
  { sequence_id: 105, metric_code: "COMMODITY.LITHIUM.CARBONATE.USD_T", name: "Lithium Carbonate", category: "COMMODITY_MACRO", value: 14350.0, unit: "USD/t", change_absolute: -100.0, change_percent: -0.69, direction: "DOWN", source: "Benchmark Minerals NA", timestamp: "2026-09-16T22:30:00Z", hash: "67ab89df" },
  { sequence_id: 106, metric_code: "MACRO.CAD_USD.SPOT_FX", name: "Bank of Canada CAD/USD", category: "COMMODITY_MACRO", value: 0.7385, unit: "CAD_per_USD", change_absolute: 0.0012, change_percent: 0.16, direction: "UP", source: "Bank of Canada Valet", timestamp: "2026-09-16T22:30:00Z", hash: "90dafe34" },
  { sequence_id: 107, metric_code: "MACRO.BOC.POLICY_RATE.PCT", name: "BoC Overnight Rate", category: "COMMODITY_MACRO", value: 2.75, unit: "percent", change_absolute: 0.0, change_percent: 0.0, direction: "FLAT", source: "Bank of Canada", timestamp: "2026-09-16T22:30:00Z", hash: "54edb821" },
  { sequence_id: 108, metric_code: "MACRO.GOC.10Y_YIELD.PCT", name: "GoC 10Y Bond Yield", category: "COMMODITY_MACRO", value: 3.02, unit: "percent", change_absolute: -0.04, change_percent: -1.31, direction: "DOWN", source: "Bank of Canada Markets", timestamp: "2026-09-16T22:30:00Z", hash: "12fa78bc" },
];

export const DEFAULT_MACRO_SUMMARY: MacroKPISummary = {
  national_emissions_abatement_mt: 42.8,
  average_indigenous_equity_pct: 18.4,
  total_ilgp_allocated_cad: 1_450_000_000,
  national_spend_run_rate_cad_mo: 285.0,
  average_domestic_content_pct: 68.2,
  average_grid_queue_wait_months: 26.4,
  average_iaac_review_duration_mo: 31.8,
  national_red_seal_deficit_fte: 14200,
  active_feed_ticks_count: 8,
  latest_commodity_ticks: INITIAL_LIVE_TICKS,
  generated_at: "2026-09-16T22:30:00Z",
};

export function evaluateProjectKPIs(project: Project): ProjectKPIScorecard {
  const regMap = new Map(CANONICAL_KPIS.map((k) => [k.code, k]));

  const assess = (code: string, observed: number, notes: string, source: string): ProjectMetricAssessment => {
    const def = regMap.get(code)!;
    const variance = observed - def.target_benchmark;
    let score = 50;
    let rank: "SUPERIOR" | "ON_TARGET" | "NEEDS_IMPROVEMENT" | "CRITICAL_GAP" = "ON_TARGET";

    const target = def.target_benchmark <= 0.0001 ? 1.0 : def.target_benchmark;

    if (def.preferred_direction === "HIGHER_BETTER") {
      const ratio = observed / target;
      score = Math.min(100, ratio * 85);
      if (ratio >= 1.15) rank = "SUPERIOR";
      else if (ratio >= 0.90) rank = "ON_TARGET";
      else if (ratio >= 0.70) rank = "NEEDS_IMPROVEMENT";
      else rank = "CRITICAL_GAP";
    } else if (def.preferred_direction === "LOWER_BETTER") {
      if (def.target_benchmark <= 0.0001) {
        if (observed <= 0.0001) {
          score = 100;
          rank = "SUPERIOR";
        } else {
          score = Math.max(10, 90 / (1 + observed * 0.1));
          if (observed > 12) rank = "CRITICAL_GAP";
          else if (observed > 3) rank = "NEEDS_IMPROVEMENT";
          else rank = "ON_TARGET";
        }
      } else {
        if (observed <= def.target_benchmark) {
          score = Math.min(100, 85 + (1 - observed / target) * 15);
          rank = "SUPERIOR";
        } else {
          score = Math.max(10, 85 / (observed / target));
          if (observed > target * 1.6) rank = "CRITICAL_GAP";
          else if (observed > target * 1.2) rank = "NEEDS_IMPROVEMENT";
          else rank = "ON_TARGET";
        }
      }
    } else {
      const dist = Math.abs(variance) / target;
      score = Math.max(20, 95 - dist * 50);
      if (dist <= 0.10) rank = "ON_TARGET";
      else if (dist <= 0.25) rank = "ON_TARGET";
      else rank = "NEEDS_IMPROVEMENT";
    }

    return {
      code: def.code,
      name: def.name,
      category: def.category,
      observed_value: Math.round(observed * 100) / 100,
      unit: def.unit,
      target_benchmark: def.target_benchmark,
      variance: Math.round(variance * 100) / 100,
      performance_rank: rank,
      score_normalized: Math.round(score * 10) / 10,
      confidence: "VERIFIED",
      notes,
      source,
    };
  };

  const isNuclear = project.sector === "Nuclear & Clean Power";
  const isMinerals = project.sector === "Critical Minerals" || project.sector === "Mining & Metals";
  const isClean = project.sector === "Clean Energy & Grid";
  const capex = project.capex_cad || 1_000_000_000;
  const capexM = capex / 1_000_000;

  // 1. ESG
  const esgMetrics: ProjectMetricAssessment[] = [
    assess("ESG.GHG.INTENSITY.SCOPE1_2", isNuclear ? 8.5 : isMinerals ? 52.0 : isClean ? 14.0 : 45.0, "Operational carbon emissions intensity", "ECCC GHG Registry"),
    assess("ESG.GHG.ABATEMENT.ANNUAL_MT", isNuclear ? 3.8 : isClean ? 2.1 : isMinerals ? 0.8 : 0.6, "Displaced counter-factual emissions", "Clean Energy Canada"),
    assess("ESG.WATER.RECYCLING.RATIO", isNuclear ? 92.0 : isMinerals ? 84.0 : 88.0, "Industrial process water closed-loop recycling", "Provincial Water Licensing"),
    assess("ESG.POWER.CLEAN_PURITY.PCT", isNuclear ? 99.5 : isClean ? 98.0 : isMinerals ? 82.0 : 85.0, "Connected non-emitting clean generation share", "System Operator Grid Balance"),
  ];

  // 2. Indigenous Equity
  const hasStrongIndig = project.name.toLowerCase().includes("indigenous") || project.name.toLowerCase().includes("oneida") || project.name.toLowerCase().includes("darlington");
  const indigEquity = hasStrongIndig ? 30.0 : isMinerals ? 15.0 : 10.0;
  const indigProcure = hasStrongIndig ? 14.0 : isMinerals ? 9.0 : 5.5;
  const indigMetrics: ProjectMetricAssessment[] = [
    assess("INDIG.EQUITY.OWNERSHIP.PCT", indigEquity, "First Nations co-investment equity stake", "CIB Indigenous Equity Database"),
    assess("INDIG.PROCUREMENT.SHARE.PCT", indigProcure, "Tier-1/Tier-2 Indigenous business contracts", "CCAB Business Directory"),
    assess("INDIG.WORKFORCE.PARTICIPATION.PCT", hasStrongIndig ? 24.0 : 16.0, "Host First Nation workforce participation", "Impact Benefit Agreement (IBA)"),
    assess("INDIG.ILGP.GUARANTEE.CAD", Math.min(150_000_000, capex * 0.10), "Secured federal loan guarantee under $5B facility", "NRCan ILGP Secretariat"),
    assess("INDIG.DIVIDEND.COMMUNITY.CAD_YR", capex * 0.11 * (indigEquity / 100) * 0.35, "30-year annualized community dividend yield", "First Nations Financial Management"),
  ];

  // 3. Capital Velocity
  const slip = project.current_stage === "DELAYED" || project.current_stage === "PAUSED" ? 14.0 : project.current_stage === "PERMITTING" ? 4.0 : 0.0;
  const capMetrics: ProjectMetricAssessment[] = [
    assess("CAPEX.SPEND.VELOCITY.CAD_MO", Math.max(8.5, capexM / 48), "Monthly capex deployment run-rate", "SEDAR+ Filings"),
    assess("SCHEDULE.SLIP.MONTHS", slip, "Months variance against FID target", "Independent Engineering Audit"),
    assess("RISK.COST_OVERRUN.HAZARD.PCT", isNuclear ? 40.6 : isMinerals ? 35.2 : 22.0, "Empirical Flyvbjerg P50 cost overrun hazard", "Oxford Megaproject Database"),
    assess("FINANCE.PUBLIC_CROWDING.MULTIPLIER", capex > 1_000_000_000 ? 4.6 : 3.8, "Private capital crowding-in multiple", "Capital Stack Structure"),
  ];

  // 4. Supply Chain
  const supplyMetrics: ProjectMetricAssessment[] = [
    assess("SUPPLY.DOMESTIC_CONTENT.PCT", isNuclear ? 84.0 : isMinerals ? 62.0 : 72.0, "Canadian BN entity procurement spend share", "Bill C-59 Verification"),
    assess("SUPPLY.BOTTLENECK.LEAD_TIME.SCORE", isNuclear ? 55.0 : isMinerals ? 48.0 : 32.0, "Critical equipment queue risk index", "Transport Canada Monitor"),
    assess("TRADE.US_TARIFF.EXPOSURE.PCT", isNuclear ? 4.0 : isMinerals ? 22.0 : 10.0, "Offtake revenue exposed to U.S. border friction", "Global Affairs Advisory"),
    assess("TRADE.IMPORT_REPLACEMENT.CAD_YR", capexM * 0.12, "Annual displacement of foreign energy or minerals", "ISED Trade Data Online"),
  ];

  // 5. Grid Physics
  const gridMetrics: ProjectMetricAssessment[] = [
    assess("GRID.PEAK_POWER.DEMAND.MW", isNuclear ? 300.0 : isMinerals ? 120.0 : 250.0, "Peak connected generation or load", "System Impact Assessment (SIA)"),
    assess("GRID.QUEUE.WAIT_TIME.MONTHS", project.province === "ON" ? 28.0 : project.province === "AB" ? 24.0 : 16.0, "Balancing authority queue latency", "IESO/AESO Interconnection Register"),
    assess("GRID.HOSTING.HEADROOM.MW", project.province === "ON" ? 85.0 : 140.0, "Available substation transformation margin", "Utility Hosting Capacity Map"),
    assess("GRID.REINFORCEMENT.COST.CAD", capex > 2_000_000_000 ? 85_000_000 : 18_000_000, "Network reinforcement upgrades expenditure", "Transmission Operator Assessment"),
  ];

  // 6. Regulatory Speed
  const isOperating = project.current_stage === "OPERATING";
  const isUnderConstruction = project.current_stage === "CONSTRUCTION";
  const regMetrics: ProjectMetricAssessment[] = [
    assess("REG.IAAC.REVIEW_DURATION.MONTHS", isOperating ? 28.0 : isUnderConstruction ? 26.0 : 22.0, "Federal Impact Assessment elapsed duration", "IAAC Registry"),
    assess("REG.PERMIT.COMPLETION.PCT", isOperating ? 100.0 : isUnderConstruction ? 95.0 : 60.0, "Statutory permits in good standing", "Provincial EA Registers"),
    assess("REG.LITIGATION.RISK_INDEX", isOperating ? 5.0 : isUnderConstruction ? 12.0 : 28.0, "Judicial review / Section 35 litigation index", "Legal Registry Analytics"),
    assess("REG.CONDITIONS.LEGALLY_BINDING.COUNT", 54.0, "Binding conditions attached to Certificate", "IAAC Decision Statement"),
  ];

  // 7. Labour Skills
  const laborMetrics: ProjectMetricAssessment[] = [
    assess("LABOR.RED_SEAL.TRADES_GAP.FTE", isNuclear ? 420.0 : 180.0, "Peak regional shortage of certified journeypersons", "BuildForce Canada Labor Model"),
    assess("LABOR.APPRENTICE.RATIO.PCT", isNuclear ? 14.5 : 12.2, "Apprentice hours worked (Bill C-59 prevailing wage bonus)", "Apprenticeship Board Audit"),
    assess("LABOR.HOUSING.ABSORPTION_DEFICIT.UNITS", isNuclear ? 140.0 : 75.0, "Local rental housing shortage within 45 min commute", "CMHC Rental Market Survey"),
    assess("LABOR.TOTAL.PERSON_YEARS.COUNT", capexM * (isNuclear ? 7.5 : 6.5), "Total lifecycle Canadian employment created", "StatCan Input-Output SUT"),
  ];

  // 8. Commodity Macro
  const macroMetrics: ProjectMetricAssessment[] = [
    assess("MACRO.BOC.POLICY_RATE.PCT", 2.75, "Bank of Canada overnight policy rate", "Bank of Canada Valet"),
    assess("MACRO.GOC.10Y_YIELD.PCT", 3.02, "GoC 10-year sovereign benchmark yield", "Bank of Canada Valet"),
    assess("MACRO.CAD_USD.SPOT_FX", 0.7385, "CAD/USD foreign exchange spot rate", "Bank of Canada Valet"),
  ];

  const calcPillar = (category: PillarScore["category"], title: string, metrics: ProjectMetricAssessment[]): PillarScore => {
    const avg = metrics.reduce((acc, m) => acc + m.score_normalized, 0) / metrics.length;
    const score = Math.round(avg * 10) / 10;
    const health = score >= 85 ? "EXEMPLARY" : score >= 70 ? "HEALTHY" : score >= 50 ? "ATTENTION_REQUIRED" : "HIGH_RISK";
    const findings = metrics
      .filter((m) => m.performance_rank === "SUPERIOR" || m.performance_rank === "CRITICAL_GAP")
      .map((m) => `${m.name} (${m.performance_rank}): ${m.observed_value} ${m.unit} vs benchmark ${m.target_benchmark} ${m.unit}`);

    return { category, title, score, health, metrics, key_findings: findings };
  };

  const pillars: PillarScore[] = [
    calcPillar("ESG_DECARBONIZATION", "ESG & Lifecycle Decarbonization", esgMetrics),
    calcPillar("INDIGENOUS_EQUITY", "Indigenous Economic Sovereignty", indigMetrics),
    calcPillar("CAPITAL_VELOCITY", "Capital Spend Velocity & Execution", capMetrics),
    calcPillar("SUPPLY_CHAIN_CONTENT", "Sovereign Supply Chain & Domestic Content", supplyMetrics),
    calcPillar("GRID_PHYSICS", "Power Grid Interconnection & Hosting", gridMetrics),
    calcPillar("REGULATORY_SPEED", "Regulatory & Statutory Permitting Latency", regMetrics),
    calcPillar("LABOR_SKILLS", "Craft Labour Supply & Apprenticeship Pipeline", laborMetrics),
    calcPillar("COMMODITY_MACRO", "Macro Benchmark & Market Offtake Alignment", macroMetrics),
  ];

  const overallRating = Math.round(
    (pillars[0].score * 0.15 +
      pillars[1].score * 0.15 +
      pillars[2].score * 0.15 +
      pillars[3].score * 0.15 +
      pillars[4].score * 0.15 +
      pillars[5].score * 0.10 +
      pillars[6].score * 0.10 +
      pillars[7].score * 0.05) *
      10
  ) / 10;

  const criticalGaps: string[] = [];
  pillars.forEach((p) => {
    p.metrics.forEach((m) => {
      if (m.performance_rank === "CRITICAL_GAP") {
        criticalGaps.push(`[${p.title}] ${m.name}: ${m.notes} (Observed: ${m.observed_value} ${m.unit} vs Target: ${m.target_benchmark} ${m.unit})`);
      }
    });
  });

  return {
    project_id: project.id,
    project_slug: project.slug,
    project_name: project.name,
    sector: project.sector,
    province: project.province,
    current_stage: project.current_stage,
    total_capex_cad: capex,
    overall_kpi_rating: overallRating,
    pillars,
    critical_action_gaps: criticalGaps,
    audit_hash: `sha256:kpi-${project.slug}-${Date.now().toString(16)}`,
    evaluated_at: new Date().toISOString(),
  };
}
