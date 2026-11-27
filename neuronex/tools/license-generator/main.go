package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/emqx/neuronex-go/internal/license"
)

func main() {
	var (
		customer   string
		contact    string
		tags       int
		expire     string
		hwToken    string
		outFile    string
		verifyCode string
		verifyFile string
		genBoth    bool
	)

	flag.StringVar(&customer, "customer", "EMQ Technologies Inc.", "Customer / Company name")
	flag.StringVar(&contact, "contact", "support@emqx.io", "Customer contact email")
	flag.IntVar(&tags, "tags", 1000000, "Maximum tag/data point capacity (e.g. 1000000 for enterprise unlimited)")
	flag.StringVar(&expire, "expire", "2999-01-01", "Expiration date (YYYY-MM-DD)")
	flag.StringVar(&hwToken, "hw", "NEURONEX-CORE-ENTERPRISE", "Hardware identifier token")
	flag.StringVar(&outFile, "out", "neuronex.lic", "Output license file path (default: neuronex.lic)")
	flag.StringVar(&verifyCode, "verify-code", "", "Verify and decode an activation code")
	flag.StringVar(&verifyFile, "verify-file", "", "Verify a .lic license file")
	flag.BoolVar(&genBoth, "all", true, "Generate both license file and activation code")
	flag.Parse()

	fmt.Println("================================================================================")
	fmt.Println("             NeuronEX License & Activation Code Generator Tool")
	fmt.Println("================================================================================")

	// Verification mode: Activation code
	if verifyCode != "" {
		fmt.Printf("[Verifying Activation Code]: %s\n\n", verifyCode)
		info, err := license.ParseActivationCode(verifyCode)
		if err != nil {
			fmt.Printf("❌ Verification FAILED: %v\n", err)
			os.Exit(1)
		}
		printLicenseSummary(info)
		fmt.Println(" Verification PASSED! Activation code is valid.")
		return
	}

	// Verification mode: License file
	if verifyFile != "" {
		fmt.Printf("[Verifying License File]: %s\n\n", verifyFile)
		data, err := os.ReadFile(verifyFile)
		if err != nil {
			fmt.Printf("❌ Failed to read file: %v\n", err)
			os.Exit(1)
		}
		var info license.LicenseInfo
		if err := json.Unmarshal(data, &info); err != nil {
			fmt.Printf("❌ Invalid JSON in license file: %v\n", err)
			os.Exit(1)
		}
		printLicenseSummary(&info)
		fmt.Println(" Verification PASSED! License file is valid.")
		return
	}

	// Normal generation mode
	fmt.Println("[Configuration]")
	fmt.Printf("  • Customer Name   : %s\n", customer)
	fmt.Printf("  • Customer Contact: %s\n", contact)
	fmt.Printf("  • Tag Capacity    : %d points\n", tags)
	fmt.Printf("  • Expire Date     : %s\n", expire)
	fmt.Printf("  • Hardware Token  : %s\n", hwToken)
	fmt.Println("--------------------------------------------------------------------------------")

	// 1. Generate Activation Code
	code := license.GenerateActivationCode(customer, tags, expire)
	fmt.Printf("\n🔑 [Generated Activation Code]:\n\n   %s\n\n", code)
	fmt.Println("   --> Copy and paste this code into NeuronEX Web UI (方式 2: 激活码) to activate!")

	// 2. Generate License File (.lic)
	expireFull := expire
	if len(expire) == 10 {
		expireFull = expire + " 00:00:00"
	}

	lic := license.LicenseInfo{
		LicenseType:     "COMMERCIAL",
		CustomerName:    customer,
		CustomerContact: contact,
		StartDate:       time.Now().Format("2006-01-02 15:04:05"),
		ExpiredDate:     expireFull,
		Expired:         false,
		TotalTagCount:   tags,
		UsedTagCount:    0,
		EnabledPlugins:  license.DefaultAllPlugins,
		HardwareToken:   hwToken,
		EkuiperOn:       true,
		OfficialEmail:   "support@emqx.io",
		Original:        false,
		Signature:       license.SignLicense(customer, tags, expire),
	}

	licData, err := json.MarshalIndent(lic, "", "  ")
	if err != nil {
		fmt.Printf("❌ Failed to encode license JSON: %v\n", err)
		os.Exit(1)
	}

	if err := os.WriteFile(outFile, licData, 0644); err != nil {
		fmt.Printf("❌ Failed to write %s: %v\n", outFile, err)
		os.Exit(1)
	}

	fmt.Printf("\n📄 [Generated License File]: %s\n", outFile)
	fmt.Println("   --> Upload this file into NeuronEX Web UI (方式 1: 上传许可证) to activate!")
	fmt.Println("================================================================================")
}

func printLicenseSummary(info *license.LicenseInfo) {
	fmt.Printf("  • Customer Name   : %s\n", info.CustomerName)
	fmt.Printf("  • Customer Contact: %s\n", info.CustomerContact)
	fmt.Printf("  • License Type    : %s\n", info.LicenseType)
	fmt.Printf("  • Total Tags Limit: %d\n", info.TotalTagCount)
	fmt.Printf("  • Expired Date    : %s\n", info.ExpiredDate)
	fmt.Printf("  • Expired         : %v\n", info.Expired)
	fmt.Printf("  • eKuiper Enabled : %v\n", info.EkuiperOn)
	fmt.Printf("  • Enabled Plugins : %d protocols enabled (%s...)\n",
		len(info.EnabledPlugins), strings.Join(info.EnabledPlugins[:min(5, len(info.EnabledPlugins))], ", "))
	fmt.Println("--------------------------------------------------------------------------------")
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
