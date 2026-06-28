import jenkins.model.Jenkins
import org.jenkinsci.plugins.workflow.job.WorkflowJob
import org.jenkinsci.plugins.workflow.cps.CpsScmFlowDefinition
import hudson.plugins.git.GitSCM
import hudson.plugins.git.UserRemoteConfig
import hudson.plugins.git.BranchSpec

// Get Jenkins instance
def jenkins = Jenkins.getInstance()

// Check if job already exists
def jobName = "uno-reverse-pipeline"
def job = jenkins.getItem(jobName)

if (job == null) {
    println("Creating pipeline job...")
    
    // Create new pipeline job
    def newJob = jenkins.createProject(WorkflowJob.class, jobName)
    newJob.setDescription("Uno Reverse Authentication System CI/CD Pipeline")
    
    // Configure Git SCM
    def gitUrl = "https://github.com/your-username/uno-reverse-auth.git"
    def userRemoteConfig = new UserRemoteConfig(gitUrl, null, null, "github-credentials")
    def scm = new GitSCM(
        [userRemoteConfig],
        [new BranchSpec("*/main")],
        false,
        [],
        null,
        null,
        []
    )
    
    // Configure pipeline
    def definition = new CpsScmFlowDefinition(scm, "Jenkinsfile")
    newJob.setDefinition(definition)
    
    // Save
    newJob.save()
    
    println("✅ Pipeline job 'uno-reverse-pipeline' created successfully!")
    println("📍 Access it at: http://localhost:8081/job/uno-reverse-pipeline/")
} else {
    println("⚠️ Job 'uno-reverse-pipeline' already exists")
}
