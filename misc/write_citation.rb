File.read('citations.txt').split("\r\n\r\n").map { |s| s.gsub('’', '\'').gsub("\r\n", "\n").gsub('—', '-') }.sort_by(&:length).each.with_index { |s, i| File.write('citations/' + i.to_s, s) }
