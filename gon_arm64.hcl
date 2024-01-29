source = ["./bin/cli-darwin-arm64"]
bundle_id = "com.kaef.cli"

apple_id {
  username = "gooogenot@gmail.com"
  provider = "PR6N2S8HV3"
}

sign {
  application_identity = "Developer ID Application: Alex Sokolov (PR6N2S8HV3)"
}

zip {
  output_path = "./bin/cli-darwin-arm64.zip"
}
