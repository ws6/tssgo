# tssgo
a golang TSS client
http://support-docs.illumina.com/SW/TSSS/TruSight_SW_API/Content/SW/FrontPages/TruSightSoftware_API.htm

Currently, found the auditlog return in very long time.
use  TSSGO_TIMEOUT_SECOND environment variable to adjust that.
or pass key TSSGO_TIMEOUT_SECOND to the NewClient  