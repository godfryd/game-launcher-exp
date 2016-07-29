games = {}
offset = 0
speed = 0
nextStop = nil
curPos = 1

function loadGamesData()
   for line in io.lines('assets/games.csv') do
      data = {}
      for w in string.gmatch(line, '([^,]+)') do
         table.insert(data, w)
      end
      game = {
         pkgName=data[1],
         name=data[2],
         screenFile=data[3],
         description=data[4]
      }
      table.insert(games, game)
   end
end

function love.load(arg)
   loadGamesData()

   for i, g in ipairs(games) do
      g.screenImg = love.graphics.newImage('assets/'..g.screenFile)
   end

   love.graphics.setNewFont("assets/SourceSansPro-Regular.ttf", 30)
end

function is_down()
   return love.keyboard.isDown('down','s')
end

function is_up()
   return love.keyboard.isDown('up','w')
end

function compute_offset(dt)
   if is_up() and curPos > 1 then
      speed = 200
      if nextStop == nil then
         nextStop = offset + 20
         curPos = curPos - 1
      end
   elseif is_down() and curPos < table.getn(games) then
      speed = -200
      if nextStop == nil then
         nextStop = offset - 20
         curPos = curPos + 1
      end
   end

   if speed ~= 0 then
      offset = offset + (speed * dt)
      if (speed > 0 and offset >= nextStop) then
         offset = nextStop
         speed = 0
         nextStop = nil
      elseif (speed < 0 and offset <= nextStop) then
         offset = nextStop
         speed = 0
         nextStop = nil
      end
   end
end

function love.update(dt)
   if love.keyboard.isDown('escape') then
      love.event.push('quit')
   end

   compute_offset(dt)
end

function love.draw(dt)
   scr_width = love.graphics.getWidth()
   scr_height = love.graphics.getHeight()
   for index, game in ipairs(games) do
      if index == curPos then
         love.graphics.setColor(255, 50, 50)
      else
         love.graphics.setColor(100, 255, 100)
      end
      pos_x = scr_width / 20
      pos_y = scr_height / 20 + (index - 1) * 30 + offset
      love.graphics.print(game.name, pos_x, pos_y)
   end

   love.graphics.setColor(255, 255, 255)
   love.graphics.draw(games[curPos].screenImg, 5 * scr_width / 20, scr_height / 20)
end
